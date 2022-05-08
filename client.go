package main

import (
	"math/rand"
	"time"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}
func (c *Client) sendLog(logs []*Log) {
	ranNum := rand.Intn(100)
	if ranNum >= 0 && ranNum < 95 {
		time.Sleep(200 * time.Millisecond)
	} else {
		time.Sleep(5000 * time.Millisecond)
	}
}
