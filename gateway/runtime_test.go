package gateway

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jacob-ebey/golang-iot/core"
)

func (peripheral *PeripheralInstance) Connect() (input <-chan core.Message, output chan core.Message, err error) {
	return peripheral.ReceiveChannel, peripheral.SendChannel, nil
}

func (peripheral *PeripheralInstance) Serializer() core.MessageSerializer {
	return &core.JsonSerializer{}
}

type shitTheBedOnInitPeripheral struct{}

var shitTheBedOnInitPeripheralError = fmt.Errorf("I shit the bed on init.")

func (*shitTheBedOnInitPeripheral) Connect() (receive <-chan core.Message, send chan core.Message, err error) {
	return nil, nil, shitTheBedOnInitPeripheralError
}

func (*shitTheBedOnInitPeripheral) Serializer() core.MessageSerializer {
	return nil
}

func TestRuntimeLogsPeripheralConnectError(t *testing.T) {
	logger := core.TestLogger{}
	runtime := &Runtime{
		Logger: &logger,
		Peripherals: []Peripheral{
			&shitTheBedOnInitPeripheral{},
		},
	}

	done := make(chan bool)
	go runtime.Execute(done)
	done <- true

	if len(logger.LoggedErrors) != 1 {
		t.Fatal("The wrong number of errors were logged.")
	}

	if logger.LoggedErrors[0] != shitTheBedOnInitPeripheralError {
		t.Fatal("The wrong error was logged.")
	}
}

type shitTheBedOnSerializeMessageSerailizer struct{}

var shitTheBedOnSerializeMessageSerailizerError = fmt.Errorf("I shit the bed on serialize.")

func (serializer *shitTheBedOnSerializeMessageSerailizer) SerialiseMessage(message core.Message) ([]byte, error) {
	return nil, shitTheBedOnSerializeMessageSerailizerError
}

func (serializer *shitTheBedOnSerializeMessageSerailizer) DeserialiseMessage(data []byte) (core.Message, error) {
	result := core.Message{}
	err := json.Unmarshal(data, &result)
	return result, err
}

func TestRuntimeForwardsPeripheralMessages(t *testing.T) {
	peripheral1ReceiveChannel := make(chan core.Message)
	peripheral1 := &PeripheralInstance{
		ReceiveChannel:    peripheral1ReceiveChannel,
		SendChannel:       make(chan core.Message),
		MessageSerializer: &core.JsonSerializer{},
	}
	peripheral2ReceiveChannel := make(chan core.Message)
	peripheral2 := &PeripheralInstance{
		ReceiveChannel:    peripheral2ReceiveChannel,
		SendChannel:       make(chan core.Message),
		MessageSerializer: &core.JsonSerializer{},
	}

	toReceive1 := core.Message{
		Type: "ROFL",
	}
	toReceive2 := core.Message{
		Type: "LOL",
	}

	sendChannel := make(chan core.Message)
	logger := core.TestLogger{}
	runtime := &Runtime{
		SendChannel:       sendChannel,
		MessageSerializer: &core.JsonSerializer{},
		Peripherals: []Peripheral{
			peripheral1,
			peripheral2,
		},
		Logger: &logger,
	}

	done := make(chan bool)
	go runtime.Execute(done)
	var written []core.Message
	go func() {
		for message := range sendChannel {
			written = append(written, message)
		}
	}()
	peripheral1ReceiveChannel <- toReceive1
	peripheral1ReceiveChannel <- toReceive1
	peripheral2ReceiveChannel <- toReceive2
	peripheral2ReceiveChannel <- toReceive2
	<-time.After(1 * time.Millisecond)
	done <- true

	for _, err := range logger.LoggedErrors {
		t.Fatal(err)
	}

	if len(written) != 4 {
		t.Fatal("Did not write the proper number of messages from peripherals.")
	}
}

func TestRuntimeForwardsPeripheralMessagesWhenConnectFailsForPeripheral(t *testing.T) {
	peripheral1ReceiveChannel := make(chan core.Message)
	peripheral1 := &PeripheralInstance{
		ReceiveChannel:    peripheral1ReceiveChannel,
		SendChannel:       make(chan core.Message),
		MessageSerializer: &core.JsonSerializer{},
	}
	peripheral2ReceiveChannel := make(chan core.Message)
	peripheral2 := &PeripheralInstance{
		ReceiveChannel:    peripheral2ReceiveChannel,
		SendChannel:       make(chan core.Message),
		MessageSerializer: &core.JsonSerializer{},
	}

	toReceive1 := core.Message{
		Type: "ROFL",
	}
	toReceive2 := core.Message{
		Type: "LOL",
	}

	sendChannel := make(chan core.Message)
	logger := core.TestLogger{}
	runtime := &Runtime{
		SendChannel:       sendChannel,
		MessageSerializer: &core.JsonSerializer{},
		Peripherals: []Peripheral{
			&shitTheBedOnInitPeripheral{},
			peripheral1,
			peripheral2,
		},
		Logger: &logger,
	}

	done := make(chan bool)
	go runtime.Execute(done)
	var written []core.Message
	go func() {
		for message := range sendChannel {
			written = append(written, message)
		}
	}()
	peripheral1ReceiveChannel <- toReceive1
	peripheral1ReceiveChannel <- toReceive1
	peripheral2ReceiveChannel <- toReceive2
	peripheral2ReceiveChannel <- toReceive2
	<-time.After(1 * time.Millisecond)
	done <- true

	if len(logger.LoggedErrors) > 1 {
		for _, err := range logger.LoggedErrors {
			t.Fatal(err)
		}
	}

	if len(written) != 4 {
		t.Fatal("Did not write the proper number of messages from peripherals.")
	}
}
