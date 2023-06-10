package main

import (
	"fmt"
	"sync"

	"github.com/90amper/metmon/internal/agent"
	"github.com/90amper/metmon/internal/config"
	"github.com/90amper/metmon/internal/logger"
)

func main() {
	agentConfig := config.AgentConfig{
		PollInterval:   fmt.Sprint(config.CmdFlags.PollInterval) + "s",   //2s
		ReportInterval: fmt.Sprint(config.CmdFlags.ReportInterval) + "s", //"10s"
		DestURL:        "http://" + config.CmdFlags.ServerUrl,
	}
	agent, err := agent.NewAgent(agentConfig)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go agent.Collector.Run(&wg)
	go agent.Sender.Run(&wg)

	logger.Log("Agent connected to " + agentConfig.DestURL)

	wg.Wait()
	fmt.Println("Service stopped")
}
