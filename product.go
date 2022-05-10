package main

import (
	uberatomic "go.uber.org/atomic"
	"sync"
)

type Product struct {
	createTimeMs     int64
	logs             []*Log
	productDataSize  *uberatomic.Int64 // the current log size in product
	productDataCount *uberatomic.Int64 // the current log count in product
	lock             sync.RWMutex
	nextSentTimeMs   int64
	result           *SentResult
	alreadySentCount int
}

func createProduct() *Product {
	nowMs := getNowTimeInMillis()
	logList := []*Log{}
	return &Product{
		createTimeMs:     nowMs,
		logs:             logList,
		productDataSize:  uberatomic.NewInt64(0),
		productDataCount: uberatomic.NewInt64(0),
		result:           createResult(),
		alreadySentCount: 0,
	}
}

func (product *Product) addLog(log *Log) {
	product.lock.Lock()
	defer product.lock.Unlock()
	product.logs = append(product.logs, log)
}
