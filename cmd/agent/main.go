package main

import (
	"fmt"
	"sync"

	"github.com/90amper/metmon/internal/agent"
	"github.com/90amper/metmon/internal/config"
)

func main() {
	agentConfig := config.AgentConfig{
		PollInterval:   "2s",
		ReportInterval: "10s",
		DestURL:        "http://localhost:8080",
	}
	agent, err := agent.NewAgent(agentConfig)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go agent.Collector.Run(&wg)
	go agent.Sender.Run(&wg)

	wg.Wait()
	fmt.Println("Service stopped")
}
