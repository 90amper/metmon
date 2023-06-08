package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/90amper/metmon/internal/collector"
	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/sender"
	"github.com/90amper/metmon/internal/storage"
)

type Agent struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	Storage        *storage.Storage
	Collector      *collector.Collector
	Sender         *sender.Sender
}

var config = models.AgentConfig{
	PollInterval:   "2s",
	ReportInterval: "10s",
	DestURL:        "http://localhost:8080",
}

func NewAgent(config models.AgentConfig) (agent *Agent, err error) {
	var a Agent
	a.PollInterval, _ = time.ParseDuration(config.PollInterval)
	a.ReportInterval, _ = time.ParseDuration(config.ReportInterval)
	a.Storage = storage.NewStorage()
	a.Collector, err = collector.NewCollector(config, a.Storage)
	if err != nil {
		return nil, err
	}
	a.Sender, err = sender.NewSender(config, a.Storage)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func main() {
	agent, err := NewAgent(config)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go agent.Collector.Run(&wg)
	go agent.Sender.Run(&wg)

	wg.Wait()
	fmt.Println("Service stopped")
	// 	work := func(id int) {
	//         defer wg.Done()
	//         fmt.Printf("Горутина %d начала выполнение \n", id)
	//         time.Sleep(2 * time.Second)
	//         fmt.Printf("Горутина %d завершила выполнение \n", id)
	//    }

	// вызываем горутины
	//    go work(1)
	//    go work(2)

	// agent.Collector.Collect()
	// agent.Collector.Collect()
	// agent.Collector.Collect()

	// // utils.PrettyPrint(agent.Storage)
	// err := agent.Sender.SendStore(*agent.Storage)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
}
