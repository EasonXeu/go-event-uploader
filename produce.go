package main

import (
	"log"
	"os"
	"time"
)

type Produce struct {
	logger         *log.Logger
	produceLogSize int64 // the current log size in produce
	logAccumulator *LogAccumulator
	broker         *Broker
	consume        *Consume

	//moverWaitGroup        *sync.WaitGroup
	//ioWorkerWaitGroup     *sync.WaitGroup
	//ioThreadPoolWaitGroup *sync.WaitGroup

	//producerConfig        *ProducerConfig

	//buckets               int
}

func InitProduce() *Produce {
	logger := log.Default()
	logger.SetOutput(os.Stdout)
	client := NewClient()
	consumeContext := initConsumeContext(logger, client)
	consume := initConsume(logger, consumeContext)
	logAccumulator := initLogAccumulator(logger)
	broker := InitBroker(logger, logAccumulator, consume)

	return &Produce{
		logger:         logger,
		produceLogSize: 0,
		logAccumulator: logAccumulator,
		broker:         broker,
		consume:        consume,
	}
}

func (produce *Produce) Start() {
	produce.logger.Println("produce start")
	//producer.moverWaitGroup.Add(1)
	//level.Info(producer.logger).Log("msg", "producer mover start")
	//go producer.mover.run(producer.moverWaitGroup, producer.producerConfig)
	//producer.ioThreadPoolWaitGroup.Add(1)
	//go producer.threadPool.start(producer.ioWorkerWaitGroup, producer.ioThreadPoolWaitGroup)
}

func getNowTimeInMillis() int64 {
	nowTimeNano := time.Now().UnixNano()
	return nowTimeNano / 1000 / 1000
}
