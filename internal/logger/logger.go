package logger

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Log(format string, args ...interface{}) {
	intformat := "%v LOG::  %s\n"
	if strings.Contains(format, "$") {
		format = strings.ReplaceAll(format, "$", "")
		intformat = "%v LOG::  %s"
	}
	// if strings.Contains(format, "^") {
	// 	format = strings.ReplaceAll(format, "^", "")
	// 	intformat = "%s"
	// }
	// for _, val := range args {
	fmt.Printf(intformat, time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
	// fmt.Printf("%v Starting server at %v\n", time.Now().Format(time.RFC3339), config.Config.ServerURL)
	// }
}

func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
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
