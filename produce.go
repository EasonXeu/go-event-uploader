package main

import (
	"errors"
	"github.com/sirupsen/logrus"
	uberatomic "go.uber.org/atomic"
	"os"
	"sync"
	"time"
)

type Produce struct {
	logger                *logrus.Logger
	produceLogSize        *uberatomic.Int64 // the current log size in produce
	logAccumulator        *LogAccumulator
	broker                *Broker
	consumerPool          *ConsumerPool
	consumerPoolWaitGroup *sync.WaitGroup
	brokerWaitGroup       *sync.WaitGroup
	producerWaitGroup     *sync.WaitGroup
	produceShutDownFlag   *uberatomic.Bool
}

func InitProduce() *Produce {
	produce := &Produce{}
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

	client := NewClient()
	consumer := initConsumer(logger, client, produce)
	retryQueue := initRetryQueue()
	consumerPool := initConsumePool(logger, consumer, retryQueue)
	producerWaitGroup := &sync.WaitGroup{}
	logAccumulator := initLogAccumulator(logger, produce, consumerPool, producerWaitGroup)
	broker := initBroker(logger, logAccumulator, consumerPool, retryQueue)

	produce.logger = logger
	produce.produceLogSize = uberatomic.NewInt64(0)
	produce.logAccumulator = logAccumulator
	produce.broker = broker
	produce.consumerPool = consumerPool
	produce.brokerWaitGroup = &sync.WaitGroup{}
	produce.consumerPoolWaitGroup = &sync.WaitGroup{}
	produce.producerWaitGroup = producerWaitGroup
	produce.produceShutDownFlag = uberatomic.NewBool(false)

	return produce
}

func (produce *Produce) Start() {
	produce.logger.Infoln("produce start")
	produce.brokerWaitGroup.Add(1)
	go produce.broker.run(produce.brokerWaitGroup)
	produce.consumerPoolWaitGroup.Add(1)
	go produce.consumerPool.start(produce.consumerPoolWaitGroup)
}

// for the sake of data not missing, we must make sure the close order follows:
// produce.produceShutDownFlag
// produce.broker.brokerShutDownFlag
// produce.consumerPool.consumerPoolShutDownFlag

func (produce *Produce) Stop() {
	produce.produceShutDownFlag.Store(true)
	produce.producerWaitGroup.Wait()
	produce.broker.brokerShutDownFlag.Store(true)
	produce.brokerWaitGroup.Wait()
	produce.consumerPool.consumerPoolShutDownFlag.Store(true)
	produce.consumerPoolWaitGroup.Wait()
	produce.logger.Infoln("produce exit")
}

func (produce *Produce) SendLog(log *Log) error {
	err := produce.waitTime()
	if err != nil {
		return err
	}
	if produce.produceShutDownFlag.Load() {
		produce.logger.Infoln("produce is shutting down, refused to send log")
		return errors.New("produce is shutting down, refused to send log")
	}
	produce.producerWaitGroup.Add(1)
	return produce.logAccumulator.addLogToProduct(log)
}

func (produce *Produce) waitTime() error {
	if MaxBlockSecondOfSendingLog > 0 {
		for i := 0; i < MaxBlockSecondOfSendingLog; i++ {
			if produce.produceLogSize.Load() > MaxSizeOfAllProducts {
				time.Sleep(time.Second)
			} else {
				return nil
			}
		}
		produce.logger.Debugln("Wait too long, exceeds the max blocking seconds")
		return errors.New(TimeoutException)
	} else if MaxBlockSecondOfSendingLog == 0 {
		if produce.produceLogSize.Load() > MaxSizeOfAllProducts {
			produce.logger.Debugln("Don't wait, exceeds the max blocking seconds")
			return errors.New(TimeoutException)
		}
		return nil
	} else {
		for {
			if produce.produceLogSize.Load() > MaxSizeOfAllProducts {
				time.Sleep(time.Second)
			} else {
				return nil
			}
		}
	}
}

func getNowTimeInMillis() int64 {
	nowTimeNano := time.Now().UnixNano()
	return nowTimeNano / 1000 / 1000
}
