package main

import (
	"container/list"
	"github.com/sirupsen/logrus"
	uberatomic "go.uber.org/atomic"
	"sync"
	"time"
)

type ConsumerPool struct {
	logger                   *logrus.Logger
	consumerPoolShutDownFlag *uberatomic.Bool
	queue                    *list.List
	lock                     sync.RWMutex
	consumer                 *Consumer
	maxConsumerCount         chan int64
	consumerWaitGroup        *sync.WaitGroup
	retryQueue               *RetryQueue
}

func initConsumePool(logger *logrus.Logger, consumer *Consumer, retryQueue *RetryQueue) *ConsumerPool {
	return &ConsumerPool{
		logger:                   logger,
		consumerPoolShutDownFlag: uberatomic.NewBool(false),
		queue:                    list.New(),
		consumer:                 consumer,
		maxConsumerCount:         make(chan int64, MaxConsumerCount),
		consumerWaitGroup:        &sync.WaitGroup{},
		retryQueue:               retryQueue,
	}
}

func (consumerPool *ConsumerPool) addTask(product *Product) {
	defer consumerPool.lock.Unlock()
	consumerPool.lock.Lock()
	consumerPool.queue.PushBack(product)
}

func (consumerPool *ConsumerPool) hasTask() bool {
	defer consumerPool.lock.Unlock()
	consumerPool.lock.Lock()
	return consumerPool.queue.Len() > 0
}

func (consumerPool *ConsumerPool) popTask() *Product {
	defer consumerPool.lock.Unlock()
	consumerPool.lock.Lock()
	if consumerPool.queue.Len() <= 0 {
		return nil
	}
	ele := consumerPool.queue.Front()
	consumerPool.queue.Remove(ele)
	return ele.Value.(*Product)
}

func (consumerPool *ConsumerPool) start(consumerPoolWaitGroup *sync.WaitGroup) {
	consumerPool.logger.Infoln("consumerPool start")
	for {
		if product := consumerPool.popTask(); product != nil {
			consumerPool.maxConsumerCount <- 1
			consumerPool.consumerWaitGroup.Add(1)
			go func(product *Product) {
				consumerPool.consumer.sendToServer(product, consumerPool.retryQueue)
				<-consumerPool.maxConsumerCount
				consumerPool.consumerWaitGroup.Done()
			}(product)
		} else {
			if consumerPool.consumerPoolShutDownFlag.Load() {
				consumerPool.logger.Infoln("consumerPool sends all products before exiting")
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	//waiting until all goroutines created by consumerPool exiting
	consumerPool.consumerWaitGroup.Wait()
	//TODO persist the data in retry queue to file
	consumerPoolWaitGroup.Done()
	consumerPool.logger.Infoln("consumerPool exit")
}
