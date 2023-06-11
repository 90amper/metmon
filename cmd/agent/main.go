package main

import (
	"fmt"
	"sync"

	"github.com/90amper/metmon/internal/agent"
	"github.com/90amper/metmon/internal/config"
	"github.com/90amper/metmon/internal/logger"
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

	logger.Log("Agent connected to " + config.Config.ServerURL)

	wg.Wait()
	fmt.Println("Service stopped")
}
