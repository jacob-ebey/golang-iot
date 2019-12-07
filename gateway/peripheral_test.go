package gateway

import (
	"context"
	"fmt"
	"testing"
)

func TestPeripheralErrorFormat(t *testing.T) {
	err := PeripheralError{
		PeripheralError: fmt.Errorf("ROFL"),
	}

	if err.Error() != "[PeripheralError] ROFL" {
		t.Fatal("message format was wrong")
	}
}

type testPeripheral struct {
	id            string
	messages      int
	errors        int
	messageToSend []byte
	errorToSend   error
}

func (peripheral *testPeripheral) ID() string {
	return peripheral.id
}

func (peripheral *testPeripheral) Listen(ctx context.Context) (chan []byte, chan error) {
	messages := make(chan []byte)
	errors := make(chan error)

	go func() {
		sentMessages := 0
		sentErrors := 0

		for i := 0; sentMessages < peripheral.messages || sentErrors < peripheral.errors; i++ {
			if sentMessages < peripheral.messages {
				messages <- peripheral.messageToSend

				sentMessages++
				continue
			}

			errors <- peripheral.errorToSend
			sentErrors++
		}
	}()

	return messages, errors
}

func (peripheral *testPeripheral) Write(ctx context.Context, payload []byte) error {
	return nil
}

func TestAzureReaderReceivesMessage(t *testing.T) {
	ctx, done := context.WithCancel(context.TODO())

	runtime := PeripheralRuntime{
		&testPeripheral{
			id:          "1",
			messages:    1,
			errors:      1,
			errorToSend: fmt.Errorf("LOL"),
		},
		&testPeripheral{
			id:          "2",
			messages:    1,
			errors:      1,
			errorToSend: fmt.Errorf("ROFL"),
		},
	}

	expectedMessages := 2
	expectedErrors := 2

	messages, errors := runtime.Listen(ctx)

	receivedMessages := []PeripheralMessage{}
	receivedPeripheralErrors := map[string]*PeripheralError{}
	unknownErrors := []error{}

	go func() {
		for {
			select {
			default:
				if len(receivedMessages) >= expectedMessages && len(receivedPeripheralErrors) >= expectedErrors {
					done()
					return
				}

			case payload := <-messages:
				receivedMessages = append(receivedMessages, payload)

			case err := <-errors:
				if peripheralErr, ok := err.(*PeripheralError); ok {
					receivedPeripheralErrors[peripheralErr.PeripheralID] = peripheralErr
					continue
				}

				unknownErrors = append(unknownErrors, err)

			case <-ctx.Done():
				return
			}

		}
	}()

	<-ctx.Done()

	if len(receivedMessages) != expectedMessages {
		t.Fatal("wrong number of messages received")
	}

	if len(receivedPeripheralErrors) != expectedErrors {
		t.Fatal("wrong number of peripheral errors received")
	}

	if len(unknownErrors) > 0 {
		t.Fatal("received unknown errors")
	}
}
