package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

func Log(format string, args ...interface{}) {
	// for _, val := range args {
	fmt.Printf("%v %s", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
	// fmt.Printf("%v Starting server at %v\n", time.Now().Format(time.RFC3339), config.Config.ServerURL)
	// }
	fmt.Println()
}

func Debug(format string, args ...interface{}) {
	// Log(format, args...)
}

func Error(err error) {
	fmt.Println(err.Error())
}

func NewDebugLogger() zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.DisableCaller = true
	lgr, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer lgr.Sync()
	return *lgr.Sugar()
}
