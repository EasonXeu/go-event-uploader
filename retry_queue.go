package main

import (
	"container/heap"
	"sync"
)

// this is a priority queue, sorted by product's nextSentTimeMs

type RetryQueue struct {
	products []*Product
	lock     sync.Mutex
}

func initRetryQueue() *RetryQueue {
	retryQueue := &RetryQueue{}
	heap.Init(retryQueue)
	return retryQueue
}

func (retryQueue *RetryQueue) addToRetryQueue(product *Product) {
	retryQueue.lock.Lock()
	defer retryQueue.lock.Unlock()
	if product != nil {
		heap.Push(retryQueue, product)
	}

}

func (retryQueue *RetryQueue) popFromRetryQueue(popAll bool) []*Product {
	nowMs := getNowTimeInMillis()
	retryQueue.lock.Lock()
	defer retryQueue.lock.Unlock()
	products := []*Product{}
	if popAll {
		for retryQueue.Len() > 0 {
			p := heap.Pop(retryQueue).(*Product)
			products = append(products, p)
		}
		return products
	}
	for retryQueue.Len() > 0 {
		p := heap.Pop(retryQueue).(*Product)
		if p.nextSentTimeMs <= nowMs {
			products = append(products, p)
		} else {
			heap.Push(retryQueue, p)
			break
		}
	}
	return products
}

// implement the interface, should not be used

func (retryQueue *RetryQueue) Push(product interface{}) {
	retryQueue.products = append(retryQueue.products, product.(*Product))
}

// implement the interface, should not be used

func (retryQueue *RetryQueue) Pop() interface{} {
	old := retryQueue.products
	n := len(old)
	x := old[n-1]
	retryQueue.products = old[0 : n-1]
	return x
}

// implement the interface

func (retryQueue *RetryQueue) Len() int {
	return len(retryQueue.products)
}

// implement the interface, should not be used

func (retryQueue *RetryQueue) Less(i, j int) bool {
	return retryQueue.products[i].nextSentTimeMs < retryQueue.products[j].nextSentTimeMs
}

// implement the interface, should not be used

func (retryQueue *RetryQueue) Swap(i, j int) {
	retryQueue.products[i], retryQueue.products[j] = retryQueue.products[j], retryQueue.products[i]
}
