package main

import "math/rand"

var (
	charSet    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321"
	charSetLen = len(charSet)
)

type Log struct {
	Content string
}

func calculateLogBytes(log *Log) int {
	return len(log.Content)
}

func calculateLogsBytes(logs []*Log) int {
	total := 0
	for _, log := range logs {
		total = total + calculateLogBytes(log)
	}
	return total
}

func generateLog() *Log {
	content := ""
	for i := 0; i < 100; i++ {
		idx := rand.Intn(charSetLen)
		content = content + string(charSet[idx])
	}
	return &Log{
		Content: content,
	}
}
