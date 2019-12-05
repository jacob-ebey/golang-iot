package gateway

import (
	"github.com/jacob-ebey/golang-iot/core"
)

type Runtime struct {
	Peripherals       []Peripheral
	SendChannel       chan core.Message
	MessageSerializer core.MessageSerializer
	Logger            core.Logger
}

func (runtime *Runtime) Execute(done chan bool) {
	peripherals := PeripheralRuntime{}

	for _, peripheral := range runtime.Peripherals {
		input, output, err := peripheral.Connect()
		if err != nil {
			runtime.Logger.LogError(err)
			continue
		}

		peripherals = append(peripherals, PeripheralInstance{
			ReceiveChannel:    input,
			SendChannel:       output,
			MessageSerializer: peripheral.Serializer(),
		})
	}

	messagesDone := make(chan bool)
	messages := peripherals.Execute(messagesDone)

	isDone := false
	for {
		select {
		case message := <-messages:
			runtime.SendChannel <- message
			break
		case isDone = <-done:
			messagesDone <- isDone
			return
		}
	}
}
