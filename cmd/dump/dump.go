package main

import (
	"fmt"
	"reflect"
	"time"

	gt7 "github.com/snipem/go-gt7-telemetry/lib"
)

func prettyPrintStruct(data interface{}) {
	val := reflect.ValueOf(data)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		fieldValue := field.Interface()

		fmt.Printf("\r%s: %+v\n", fieldName, fieldValue)
	}
}

func main() {
	gt7c := gt7.NewGT7Communication("192.168.1.215")
	go gt7c.Run()
	for true {
		fmt.Print("\033[H\033[2J")
		prettyPrintStruct(gt7c.LastData)
		time.Sleep(60 * time.Millisecond)
	}
}
