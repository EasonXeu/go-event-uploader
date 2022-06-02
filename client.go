package main

import (
	"errors"
	uberatomic "go.uber.org/atomic"
	"math/rand"
	"time"
)

type Client struct {
}

var SendCount = uberatomic.NewInt64(0)
var ErrorCount = uberatomic.NewInt64(0)

func NewClient() *Client {
	return &Client{}
}
func (c *Client) sendLog(product *Product) error {
	// you can compress your logs before sending them by the network
	// I recommend to use lz4, you can find some benchmarks here
	// (https://catchchallenger.first-world.info/wiki/Quick_Benchmark:_Gzip_vs_Bzip2_vs_LZMA_vs_XZ_vs_LZ4_vs_LZO)
	logs := product.logs
	ranNum := rand.Intn(10000)
	if ranNum >= 0 && ranNum < 9000 {
		time.Sleep(200 * time.Millisecond)
		SendCount.Add(int64(len(logs)))
		return nil
	} else if ranNum >= 9000 && ranNum < 9900 {
		time.Sleep(5000 * time.Millisecond)
		SendCount.Add(int64(len(logs)))
		return nil
	} else {
		ErrorCount.Add(int64(len(logs)))
		return errors.New("send error")
	}
}
