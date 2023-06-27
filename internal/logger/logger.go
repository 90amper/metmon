package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func Log(args ...interface{}) {
	for _, val := range args {
		fmt.Printf("%+v\t", val)
	}
}

func NewDebugLogger() zap.SugaredLogger {
	// создаём предустановленный регистратор zap
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.DisableCaller = true
	lgr, err := config.Build()
	// lgr, err := zap.NewProduction()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer lgr.Sync()
	// делаем регистратор SugaredLogger
	return *lgr.Sugar()
}
