package logger

import "fmt"

func Log(args ...interface{}) {
	for _, val := range args {
		fmt.Printf("%+v\n", val)
	}
}
