package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/90amper/metmon/internal/agent"
	"github.com/90amper/metmon/internal/agent/config"
)

func main() {
	agent, err := agent.NewAgent()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go agent.Collector.Run(&wg)
	go agent.Sender.Run(&wg)
	fmt.Printf("%v Starting agent, connecting to %v\n", time.Now().Format(time.RFC3339), config.Config.ServerURL)

	wg.Wait()
	fmt.Printf("%v Agent stopped\n", time.Now().Format(time.RFC3339))
}
