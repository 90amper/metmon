package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func Log(args ...interface{}) {
	for _, val := range args {
		fmt.Printf("%+v\t", val)
	}
	fmt.Println()
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
