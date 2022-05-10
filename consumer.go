package main

import (
	"github.com/sirupsen/logrus"
	"math"
)

type Consumer struct {
	logger  *logrus.Logger
	client  *Client
	produce *Produce
}

func initConsumer(logger *logrus.Logger, client *Client, produce *Produce) *Consumer {
	return &Consumer{
		logger:  logger,
		client:  client,
		produce: produce,
	}
}

func (consumer *Consumer) sendToServer(product *Product, retryQueue *RetryQueue) {
	beginMs := getNowTimeInMillis()
	err := consumer.client.sendLog(product)
	endMs := getNowTimeInMillis()
	costMs := endMs - beginMs
	if err == nil {
		consumer.logger.Debugln("customer succeed to send log to server")
		resultDetail := createResultDetail(true, "", "", "", "", beginMs, costMs)
		consumer.addResultToProduct(product, true, resultDetail)
		// we should release the space occupied by the data sent where successfully
		consumer.produce.produceLogSize.Sub(product.productDataSize.Load())
		return
	}

	resultDetail := createResultDetail(false, "requestId", "httpCode", "errorCode", "errorMsg", beginMs, costMs)
	consumer.addResultToProduct(product, false, resultDetail)
	if consumer.doNotNeedToRetry(err) {
		consumer.logger.Debugln("customer failed to send log, do not need to retry so abandon it.")
		// we should release the space occupied by the data sent before abandoning
		consumer.produce.produceLogSize.Sub(product.productDataSize.Load())
		return
	}
	if product.alreadySentCount >= 1+MaxRetryCount {
		consumer.logger.Debugln("customer failed to send log, exceeds the max retry count so abandon it.")
		// we should release the space occupied by the data sent before abandoning
		consumer.produce.produceLogSize.Sub(product.productDataSize.Load())
		return
	}
	consumer.logger.Debugln("customer failed to send log, adds to retry queue.")
	retryBackoffTime := math.Min(float64(MaxRetryBackoffMs), float64(BaseRetryBackoffMs)*math.Pow(2, float64(product.alreadySentCount)))
	product.nextSentTimeMs = getNowTimeInMillis() + int64(retryBackoffTime)
	retryQueue.addToRetryQueue(product)
	// we should not release the space occupied by the data sent because the data is neither being sent nor being abandoned
	return
}

func (consumer *Consumer) doNotNeedToRetry(e error) bool {
	// implements your errors which are not retried
	return false
}

func (consumer *Consumer) addResultToProduct(product *Product, success bool, resultDetail *ResultDetail) {
	product.result.addResultDetail(resultDetail)
	product.alreadySentCount += 1
	product.result.success = success
}
