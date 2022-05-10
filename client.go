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
