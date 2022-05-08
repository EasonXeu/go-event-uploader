package main

import (
	uberatomic "go.uber.org/atomic"
	"log"
	"time"
)

type Broker struct {
	logger             *log.Logger
	brokerShutDownFlag *uberatomic.Bool

	logAccumulator *LogAccumulator //data source
	consume        *Consume        //data destination

	//retryQueue        *RetryQueue
	//ioWorker          *IoWorker
}

func InitBroker(logger *log.Logger, logAccumulator *LogAccumulator, consume *Consume) *Broker {
	return &Broker{
		logger:             logger,
		brokerShutDownFlag: uberatomic.NewBool(false),
		logAccumulator:     logAccumulator,
		consume:            consume,
	}
}

func (broker *Broker) run() {
	for !broker.brokerShutDownFlag.Load() {
		sleepMs := MaxWaitingMsOfProduct
		nowTimeMs := getNowTimeInMillis()
		broker.logAccumulator.lock.Lock()
		for key, product := range broker.logAccumulator.logData {
			timeBeforeDeadline := product.createTimeMs + MaxWaitingMsOfProduct - nowTimeMs
			if timeBeforeDeadline <= 0 {
				broker.logger.Println("broker goroutine moves product to consume")
				broker.consume.addTask(product)
				delete(broker.logAccumulator.logData, key)
			} else {
				if sleepMs > timeBeforeDeadline {
					sleepMs = timeBeforeDeadline
				}
			}
		}
		broker.logAccumulator.lock.Unlock()
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)

		//retryProducerBatchList := mover.retryQueue.getRetryBatch(mover.moverShutDownFlag.Load())
		//if retryProducerBatchList == nil {
		// If there is nothing to send in the retry queue, just wait for the minimum time that was given to me last time.
		//time.Sleep(time.Duration(sleepMs) * time.Millisecond)
		//} else {
		//	count := len(retryProducerBatchList)
		//	for i := 0; i < count; i++ {
		//		mover.threadPool.addTask(retryProducerBatchList[i])
		//	}
		//}
	}

}
