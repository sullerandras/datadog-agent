package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

func main() {
	fmt.Println("waiting a few")
	time.Sleep(5 * time.Second)
	fmt.Println("starting flood")
	c, err := statsd.New("dogstatsd:8125")
	if err != nil {
		fmt.Printf("could not connect to dogstatsd: %s\n", err)
		return
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Starting client")
	for {
		select {
		// Block here until we receive the interrupt signal
		case <-signalCh:
			fmt.Println("Stopping client")
			return
		default:
			if err := c.Incr("bench_counter", nil, 1); err != nil {
				fmt.Printf("Could not send package: %s\n", err)
				return
			}
		}
	}
}
