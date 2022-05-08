package main

const (
	MaxWaitingMsOfProduct int64 = 2 * 1000              // the max waiting time before being sent
	MaxSizeOfProduct            = 512                   // the max bytes size before being sent (KB)
	MaxCountOfProduct           = 1000                  // the max count before being sent (KB)
	MaxSizeInTotal              = MaxSizeOfProduct * 20 // the max byte size of produce (KB)
)
