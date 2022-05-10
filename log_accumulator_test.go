package main

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestHashKey(t *testing.T) {
	var hashSize int64
	hashSize = 20
	lock := sync.Mutex{}
	statistics := make(map[string]int64)
	wg := sync.WaitGroup{}
	concurrency := 10
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			for i := 0; i < 20000; i++ {
				hashK := hashKey(nil, hashSize)
				lock.Lock()
				val, exist := statistics[hashK]
				if !exist {
					statistics[hashK] = 1
				} else {
					statistics[hashK] = val + 1
				}
				lock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, hashSize, int64(len(statistics)))
}
