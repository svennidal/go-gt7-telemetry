package main

import (
	"fmt"
	"os/exec"
	"reflect"
	"time"

	gt7 "github.com/snipem/go-gt7-telemetry/lib"
)

const (
	siri = "say"
)

func talk(message string) {
	cmd := exec.Command(siri, message)

	_, err := cmd.Output()
	if err != nil {
		panic(err.Error())
	}
}

func remainingLaps(total, current int16) string {
	remaining := (total - current) + 1

	if remaining < 0 || total < 0 || current < 0 || remaining > total {
		return ""
	} else if remaining == 1 {
		return "1 lap left! Come on!"
	} else if current > total {
		return "Alright Cowboy! You made it!"
	} else {
		return fmt.Sprintf("%d laps left!", remaining)
	}

	return ""
}

func remainingFuel(fuelCapacity, currentFuel float32, totalLaps, currentLap int16) string {
	if currentFuel == 100 {
		return ""
	}

	spent := fuelCapacity - currentFuel
	perRing := spent / float32(currentLap)
	remaining := int(currentFuel / perRing)

	message := fmt.Sprintf("%d laps left of gas.", remaining)

	if remaining < 2 {
		message = message + " You need to gas up now!"
	}

	return message
}

func changeGear(suggested, current uint8) string {
	move := int16(suggested) - int16(current)

	message := fmt.Sprintf("next turn %d... you're at %d.", suggested, current)

	if move < 0 {
		move *= -1
		message = fmt.Sprintf("%s Move %d down!", message, move)
	} else if move > 0 {
		message = fmt.Sprintf("%s Move %d up!", message, move)
	} else {
		message = fmt.Sprintf("%s You're good!", message)
	}

	return message
}

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

type previous struct {
	CurrentLap           int16
	CurrentSuggestedGear uint8
}

func main() {
	pre := previous{}

	talk("connecting")
	gt7c := gt7.NewGT7Communication("192.168.1.215")
	go gt7c.Run()
	fmt.Printf("connected and running.\n")

	pre.CurrentLap = gt7c.LastData.CurrentLap

	for true {
		//fmt.Print("\033[H\033[2J")
		//prettyPrintStruct(gt7c.LastData)

		if pre.CurrentLap != gt7c.LastData.CurrentLap {
			lapMessage := remainingLaps(
				gt7c.LastData.TotalLaps,
				gt7c.LastData.CurrentLap,
			)
			talk(lapMessage)

			fuelMessage := remainingFuel(
				gt7c.LastData.FuelCapacity,
				gt7c.LastData.CurrentFuel,
				gt7c.LastData.TotalLaps,
				gt7c.LastData.CurrentLap,
			)
			talk(fuelMessage)

			pre.CurrentLap = gt7c.LastData.CurrentLap
		}

		if gt7c.LastData.SuggestedGear != 15 {
			if pre.CurrentSuggestedGear != gt7c.LastData.SuggestedGear {
				gearMessage := changeGear(gt7c.LastData.SuggestedGear, gt7c.LastData.CurrentGear)
				talk(gearMessage)
			}
		}

		pre.CurrentSuggestedGear = gt7c.LastData.SuggestedGear

		time.Sleep(100 * time.Millisecond)
	}
}
