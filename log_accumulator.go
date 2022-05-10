package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	uberatomic "go.uber.org/atomic"
	"math/rand"
	"sync"
)

const hashKeyFormat = "h-%d"

type LogAccumulator struct {
	logger                     *logrus.Logger
	logAccumulatorShutDownFlag *uberatomic.Bool
	lock                       sync.Mutex // to protect logData
	logData                    map[string]*Product
	logDataMapHashSize         int64
	produce                    *Produce
	consume                    *ConsumerPool
	producerWaitGroup          *sync.WaitGroup
}

func initLogAccumulator(logger *logrus.Logger, produce *Produce, consume *ConsumerPool, producerWaitGroup *sync.WaitGroup) *LogAccumulator {
	hashSize := MaxSizeOfAllProducts / MaxSizeOfProduct
	return &LogAccumulator{
		logger:                     logger,
		logAccumulatorShutDownFlag: uberatomic.NewBool(false),
		logData:                    make(map[string]*Product, hashSize),
		logDataMapHashSize:         hashSize,
		produce:                    produce,
		consume:                    consume,
		producerWaitGroup:          producerWaitGroup,
	}
}

func (logAccumulator *LogAccumulator) addLogToProduct(log *Log) error {
	hashIndex := hashKey(log, logAccumulator.logDataMapHashSize)
	logSize := int64(calculateLogBytes(log))
	// put log to map, thread-safe
	logAccumulator.lock.Lock()
	if _, exist := logAccumulator.logData[hashIndex]; !exist {
		createdProduct := createProduct()
		logAccumulator.logData[hashIndex] = createdProduct
	}
	destProduct := logAccumulator.logData[hashIndex]
	destProduct.addLog(log)
	destProduct.productDataSize.Add(logSize)
	destProduct.productDataCount.Add(1)
	logAccumulator.produce.produceLogSize.Add(logSize)
	logAccumulator.moveToConsumeIfNecessary(hashIndex, destProduct)
	logAccumulator.lock.Unlock()
	logAccumulator.producerWaitGroup.Done()
	return nil
}

func (logAccumulator *LogAccumulator) moveToConsumeIfNecessary(key string, product *Product) {
	if product.productDataSize.Load() >= MaxSizeOfProduct || product.productDataCount.Load() >= MaxCountOfProduct {
		logAccumulator.logger.Debugln("producer moves product to consumerPool")
		logAccumulator.consume.addTask(product)
		delete(logAccumulator.logData, key)
	}
}

func hashKey(log *Log, hashSize int64) string {
	//implement yourself hash here
	ranNm := rand.Int63n(hashSize)
	return fmt.Sprintf(hashKeyFormat, ranNm)
}
