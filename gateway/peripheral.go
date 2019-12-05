package gateway

import (
	"github.com/jacob-ebey/golang-iot/core"
)

type Peripheral interface {
	Connect() (receive <-chan core.Message, send chan core.Message, err error)
	Serializer() core.MessageSerializer
}

type PeripheralInstance struct {
	ReceiveChannel    <-chan core.Message
	SendChannel       chan core.Message
	MessageSerializer core.MessageSerializer
}

type PeripheralRuntime []PeripheralInstance

func (peripherals PeripheralRuntime) Execute(done chan bool) (messages chan core.Message) {
	messages = make(chan core.Message)

	for _, peripheral := range peripherals {
		go func(peripheral PeripheralInstance) {
			for {
				select {
				case message := <-peripheral.ReceiveChannel:
					messages <- message
				case done := <-done:
					if done {
						return
					}
				}
			}
		}(peripheral)
	}

	return messages
}
