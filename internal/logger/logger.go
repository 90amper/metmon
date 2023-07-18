package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Log(format string, args ...interface{}) {
	// for _, val := range args {
	fmt.Printf("%v LOG::  %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
	// fmt.Printf("%v Starting server at %v\n", time.Now().Format(time.RFC3339), config.Config.ServerURL)
	// }
}

func Debug(format string, args ...interface{}) {
	fmt.Printf("%v DEBUG::  %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
}

func Trace(format string, args ...interface{}) {
	// fmt.Printf("%v TRACE::  %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
}

func Error(err error) {
	fmt.Printf("%v ERROR::  %s\n", time.Now().Format(time.RFC3339), err.Error())
}

func NewDebugLogger() zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.OutputPaths = []string{"stdout"}
	config.DisableCaller = true
	lgr, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer lgr.Sync()
	return *lgr.Sugar()
}
