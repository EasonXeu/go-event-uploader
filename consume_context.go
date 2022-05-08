package main

import "log"

type ConsumeContext struct {
	logger *log.Logger
	client *Client
	//taskCount              int64
	//retryQueue             *RetryQueue
	//retryQueueShutDownFlag *uberatomic.Bool
	//maxIoWorker            chan int64
	//noRetryStatusCodeMap   map[int]*string
	//producer               *Producer
}

func initConsumeContext(logger *log.Logger, client *Client) *ConsumeContext {
	return &ConsumeContext{
		logger: logger,
		client: client,
	}
}
