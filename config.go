package main

const (
	MaxWaitingMsOfProduct      int64 = 2 * 1000               // the max waiting time before being sent, default 2000ms
	MaxSizeOfProduct           int64 = 512 * 1024             // the max bytes size of product (B), default 512KB. Products whose size are larger than this will be sent
	MaxCountOfProduct          int64 = 1000                   // the max count of product, default 100. Products whose count are more than this will be sent
	MaxSizeOfAllProducts       int64 = 100 * MaxSizeOfProduct // the max byte size of all product (B), default 50MB. The memory size occupied by all products
	MaxBlockSecondOfSendingLog int   = 30                     // the max block seconds if the produce available space is not enough, default 30s. The block time for what invokers are willing to wait
	MaxConsumerCount           int   = 50                     // the max count of goroutine, default 50. The count of goroutines which are used to send log
	MaxRetryCount              int   = 10                     // the max retry count before abandoning it, default 10. The logs won't be abandoned until this value is reached
	BaseRetryBackoffMs         int64 = 100                    // the base retry backoff time, default 100ms. The base back off time
	MaxRetryBackoffMs          int64 = 50 * 1000              // the max retry backoff time, default 50s. The max back off time
)

type Config struct {
	maxWaitingMsOfProduct      int64
	maxSizeOfProduct           int64
	maxCountOfProduct          int64
	maxSizeOfAllProducts       int64
	maxBlockSecondOfSendingLog int
	maxConsumerCount           int
}

type ConfigBuilder struct {
	maxWaitingMsOfProduct      int64
	maxSizeOfProduct           int64
	maxCountOfProduct          int64
	maxSizeOfAllProducts       int64
	maxBlockSecondOfSendingLog int
	maxConsumerCount           int
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		maxWaitingMsOfProduct:      MaxWaitingMsOfProduct,
		maxSizeOfProduct:           MaxSizeOfProduct,
		maxCountOfProduct:          MaxCountOfProduct,
		maxSizeOfAllProducts:       MaxSizeOfAllProducts,
		maxBlockSecondOfSendingLog: MaxBlockSecondOfSendingLog,
		maxConsumerCount:           MaxConsumerCount,
	}
}

func (builder *ConfigBuilder) SetMaxWaitingMsOfProduct(val int64) *ConfigBuilder {
	if val > 500 {
		builder.maxWaitingMsOfProduct = val
	}
	return builder
}

func (builder *ConfigBuilder) SetMaxSizeOfProduct(val int64) *ConfigBuilder {
	if val > 0 && val <= 5*1024*1024 {
		builder.maxSizeOfProduct = val
	}
	return builder
}

func (builder *ConfigBuilder) SetMaxCountOfProduct(val int64) *ConfigBuilder {
	if val > 0 && val <= 40960 {
		builder.maxCountOfProduct = val
	}
	return builder
}

func (builder *ConfigBuilder) SetMaxSizeOfAllProducts(val int64) *ConfigBuilder {
	if val > 0 && val <= 100*1024*1024 {
		builder.maxSizeOfAllProducts = val
	}
	return builder
}

func (builder *ConfigBuilder) SetMaxBlockSecondOfSendingLog(val int) *ConfigBuilder {
	builder.maxBlockSecondOfSendingLog = val
	return builder
}

func (builder *ConfigBuilder) SetMaxConsumerCount(val int) *ConfigBuilder {
	if val > 0 {
		builder.maxConsumerCount = val
	}
	return builder
}

func (builder *ConfigBuilder) Build() *Config {
	return &Config{
		maxWaitingMsOfProduct:      builder.maxWaitingMsOfProduct,
		maxSizeOfProduct:           builder.maxSizeOfProduct,
		maxCountOfProduct:          builder.maxCountOfProduct,
		maxSizeOfAllProducts:       builder.maxSizeOfAllProducts,
		maxBlockSecondOfSendingLog: builder.maxBlockSecondOfSendingLog,
		maxConsumerCount:           builder.maxConsumerCount,
	}
}
