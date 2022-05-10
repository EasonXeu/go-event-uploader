package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRetryQueue_NotAll(t *testing.T) {
	nowMs := getNowTimeInMillis()
	p1 := &Product{
		nextSentTimeMs: nowMs + 2000,
	}
	p2 := &Product{
		nextSentTimeMs: nowMs + 4000,
	}
	p3 := &Product{
		nextSentTimeMs: nowMs + 6000,
	}
	p4 := &Product{
		nextSentTimeMs: nowMs + 8000,
	}
	p5 := &Product{
		nextSentTimeMs: nowMs + 10000,
	}
	p6 := &Product{
		nextSentTimeMs: nowMs + 9000,
	}
	retryQueue := initRetryQueue()
	retryQueue.addToRetryQueue(p2)
	retryQueue.addToRetryQueue(p1)
	retryQueue.addToRetryQueue(p4)
	retryQueue.addToRetryQueue(p3)
	retryQueue.addToRetryQueue(p5)

	pList := retryQueue.popFromRetryQueue(false)
	assert.True(t, 0 == len(pList))

	time.Sleep(2000 * time.Millisecond)
	pList = retryQueue.popFromRetryQueue(false)
	assert.True(t, 1 == len(pList))
	assert.True(t, nowMs+2000 == pList[0].nextSentTimeMs)

	time.Sleep(1000 * time.Millisecond)
	pList = retryQueue.popFromRetryQueue(false)
	assert.True(t, 0 == len(pList))

	time.Sleep(1000 * time.Millisecond)
	pList = retryQueue.popFromRetryQueue(false)
	assert.True(t, 1 == len(pList))
	assert.True(t, nowMs+4000 == pList[0].nextSentTimeMs)

	time.Sleep(4000 * time.Millisecond)
	pList = retryQueue.popFromRetryQueue(false)
	assert.True(t, 2 == len(pList))
	assert.True(t, nowMs+6000 == pList[0].nextSentTimeMs)
	assert.True(t, nowMs+8000 == pList[1].nextSentTimeMs)

	retryQueue.addToRetryQueue(p6)
	pList = retryQueue.popFromRetryQueue(false)
	assert.True(t, 0 == len(pList))

	time.Sleep(1000 * time.Millisecond)
	pList = retryQueue.popFromRetryQueue(false)
	assert.True(t, 1 == len(pList))
	assert.True(t, nowMs+9000 == pList[0].nextSentTimeMs)

	time.Sleep(500 * time.Millisecond)
	pList = retryQueue.popFromRetryQueue(false)
	assert.True(t, 0 == len(pList))

	time.Sleep(500 * time.Millisecond)
	pList = retryQueue.popFromRetryQueue(false)
	assert.True(t, 1 == len(pList))
	assert.True(t, nowMs+10000 == pList[0].nextSentTimeMs)

}

func TestRetryQueue_All(t *testing.T) {
	nowMs := getNowTimeInMillis()
	p1 := &Product{
		nextSentTimeMs: nowMs + 2000,
	}
	p2 := &Product{
		nextSentTimeMs: nowMs + 4000,
	}
	p3 := &Product{
		nextSentTimeMs: nowMs + 6000,
	}
	p4 := &Product{
		nextSentTimeMs: nowMs + 8000,
	}
	p5 := &Product{
		nextSentTimeMs: nowMs + 10000,
	}

	retryQueue := initRetryQueue()
	retryQueue.addToRetryQueue(p2)
	retryQueue.addToRetryQueue(p1)
	retryQueue.addToRetryQueue(p4)
	retryQueue.addToRetryQueue(p3)
	retryQueue.addToRetryQueue(p5)

	pList := retryQueue.popFromRetryQueue(true)
	assert.True(t, 5 == len(pList))
	assert.True(t, nowMs+2000 == pList[0].nextSentTimeMs)
	assert.True(t, nowMs+4000 == pList[1].nextSentTimeMs)
	assert.True(t, nowMs+6000 == pList[2].nextSentTimeMs)
	assert.True(t, nowMs+8000 == pList[3].nextSentTimeMs)
	assert.True(t, nowMs+10000 == pList[4].nextSentTimeMs)

}
