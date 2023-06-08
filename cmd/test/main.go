package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
)

func PrettyPrint(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return nil
}

func main() {
	runtMetrics := GetRuntimeMetrics()
	PrettyPrint(runtMetrics)
}

func GetRuntimeMetrics() map[string]float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// PrettyPrint(m)
	values := reflect.ValueOf(m)
	numfield := values.NumField()
	// f64c := 0
	res := make(map[string]float64)
	for i := 0; i < numfield; i++ {
		// field := values.Field(i)
		fName := reflect.TypeOf(m).Field(i).Name
		fType := reflect.ValueOf(m).Field(i).Type().String()
		fValue := reflect.ValueOf(m).Field(i)
		if fType == "uint64" || fType == "uint32" {
			// f64c++
			// fmt.Printf("%20v\t%-10v\t%v\n", fName, fType, fValue.Uint())
			res[fName] = float64(fValue.Uint())
		} else if fType == "float64" {
			// fmt.Printf("%20v\t%-10v\t%v\n", fName, fType, fValue.Float())
			res[fName] = fValue.Float()
		} else {
			fmt.Printf("skip %100v\t%v\n", fName, fType)
		}
	}
	return res
	// println(f64c)
	// for k, v := range m {
	// 	fmt.Printf("(%T)\t%v\t%v", v, k, v)
	// }
}
