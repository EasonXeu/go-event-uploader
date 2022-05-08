package main

import (
	"container/list"
	uberatomic "go.uber.org/atomic"
	"log"
	"sync"
)

type Consume struct {
	logger              *log.Logger
	consumeShutDownFlag *uberatomic.Bool

	queue          *list.List
	lock           sync.RWMutex
	consumeContext *ConsumeContext
}

func initConsume(logger *log.Logger, consumeContext *ConsumeContext) *Consume {
	return &Consume{
		logger:              logger,
		consumeShutDownFlag: uberatomic.NewBool(false),
		queue:               list.New(),
		consumeContext:      consumeContext,
	}
}

func (consume *Consume) addTask(product *Product) {
	defer consume.lock.Unlock()
	consume.lock.Lock()
	consume.queue.PushBack(product)
}

func (consume *Consume) hasTask() bool {
	defer consume.lock.Unlock()
	consume.lock.Lock()
	return consume.queue.Len() > 0
}

func (consume *Consume) popTask() *Product {
	defer consume.lock.Unlock()
	consume.lock.Lock()
	if consume.queue.Len() <= 0 {
		return nil
	}
	ele := consume.queue.Front()
	consume.queue.Remove(ele)
	return ele.Value.(*Product)
}
