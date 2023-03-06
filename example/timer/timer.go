package main

import (
	"log"
	"time"
)

var (
	Plugin      Timer
	Name        = "Timer"
	Description = "a simple plugin, print 'time ticker' every second"
	Version     = "1.0.0"
)

type (
	Timer struct{}
)

func (t Timer) Run() {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		log.Println("time ticker")
	}
}
