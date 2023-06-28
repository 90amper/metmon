package main

import (
	"github.com/90amper/metmon/internal/server"
)

func main() {
	// var signals = []os.Signal{
	// 	os.Interrupt,
	// 	syscall.SIGINT,
	// 	syscall.SIGQUIT,
	// 	syscall.SIGABRT,
	// 	syscall.SIGKILL,
	// 	syscall.SIGTERM,
	// 	// syscall.SIGSTOP,
	// }
	// shutdown := make(chan os.Signal, 1)
	// signal.Notify(shutdown, signals...)
	// rootCtx := context.Background()
	// taskCtx, cancelFn := context.WithCancel(rootCtx)
	server.Run()

	// <-shutdown
	// cancelFn()
	// fmt.Printf("Shutdown MetMon server at %v", time.Now().Format(time.RFC3339))
}
