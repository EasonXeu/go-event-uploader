package main

import "sync"

type Product struct {
	createTimeMs  int64
	logs          []*Log
	lock          sync.RWMutex
	totalDataSize int64
	result        *Result

	//maxRetryIntervalInMs int64
	//maxRetryTimes        int
	//baseRetryBackoffMs   int64
	//maxReservedAttempts  int

	//attemptCount         int  //already tried
	//nextRetryMs          int64

	//logsSize         int
	//logsCount        int
	//callBackList         []CallBack
	//project              string
	//logstore             string
	//shardHash            *string
}
