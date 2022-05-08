package main

import (
	uberatomic "go.uber.org/atomic"
	"log"
	"sync"
)

const defaultProductKeySuffix = "log-"

type LogAccumulator struct {
	logger                     *log.Logger
	logAccumulatorShutDownFlag uberatomic.Bool
	lock                       sync.Mutex
	logData                    map[string]*Product
}

func initLogAccumulator(logger *log.Logger) *LogAccumulator {
	return &LogAccumulator{
		logData:                    make(map[string]*Product),
		logAccumulatorShutDownFlag: *uberatomic.NewBool(false),
		logger:                     logger,
	}
}
