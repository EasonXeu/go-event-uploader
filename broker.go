package main

import (
	"github.com/sirupsen/logrus"
	uberatomic "go.uber.org/atomic"
	"math"
	"sync"
	"time"
)

type Broker struct {
	logger             *logrus.Logger
	brokerShutDownFlag *uberatomic.Bool
	logAccumulator     *LogAccumulator //data source
	consume            *ConsumerPool   //data destination
	retryQueue         *RetryQueue     //data source
}

const defaultSleepMs float64 = 2000

func initBroker(logger *logrus.Logger, logAccumulator *LogAccumulator, consume *ConsumerPool, retryQueue *RetryQueue) *Broker {
	return &Broker{
		logger:             logger,
		brokerShutDownFlag: uberatomic.NewBool(false),
		logAccumulator:     logAccumulator,
		consume:            consume,
		retryQueue:         retryQueue,
	}
}

func (broker *Broker) run(brokerWaitGroup *sync.WaitGroup) {
	broker.logger.Infoln("broker start")
	for !broker.brokerShutDownFlag.Load() {
		broker.logger.Debugln("broker run once")
		var sleepMs int64 = int64(math.Min(defaultSleepMs, float64(MaxWaitingMsOfProduct)))
		nowTimeMs := getNowTimeInMillis()
		broker.logAccumulator.lock.Lock()
		for key, product := range broker.logAccumulator.logData {
			timeBeforeDeadline := product.createTimeMs + MaxWaitingMsOfProduct - nowTimeMs
			if timeBeforeDeadline <= 0 {
				broker.logger.Debugln("broker moves product from logAccumulator to consumerPool")
				broker.consume.addTask(product)
				delete(broker.logAccumulator.logData, key)
			} else {
				if sleepMs > timeBeforeDeadline {
					sleepMs = timeBeforeDeadline
				}
			}
		}
		broker.logAccumulator.lock.Unlock()

		productsToSentAgain := broker.retryQueue.popFromRetryQueue(false)
		if productsToSentAgain != nil && len(productsToSentAgain) > 0 {
			broker.logger.Debugln("broker moves product from retryQueue to consumerPool")
			for _, p := range productsToSentAgain {
				broker.consume.addTask(p)
			}
		}

		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	}
	broker.logger.Infoln("broker moves all products from logAccumulator to consumerPool before exiting")
	broker.logAccumulator.lock.Lock()
	for key, product := range broker.logAccumulator.logData {
		broker.consume.addTask(product)
		delete(broker.logAccumulator.logData, key)
	}
	broker.logAccumulator.logData = make(map[string]*Product) // release the memory
	broker.logAccumulator.lock.Unlock()

	broker.logger.Infoln("broker moves all products from retryQueue to consumerPool before exiting")
	productsToRetry := broker.retryQueue.popFromRetryQueue(true)
	for _, p := range productsToRetry {
		broker.consume.addTask(p)
	}
	brokerWaitGroup.Done()
	broker.logger.Infoln("broker exit")
}
