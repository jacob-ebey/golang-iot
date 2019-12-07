package gateway

import (
	"context"
	"fmt"
)

// Peripheral describes a peripheral device.
type Peripheral interface {
	ID() string
	Listen(ctx context.Context) (chan []byte, <-chan error)
	Write(ctx context.Context, payload []byte) error
}

// PeripheralRuntime is a list of Peripheral devices that expose functionality.
type PeripheralRuntime []Peripheral

// PeripheralMessage describes a message received from a peripheral device.
type PeripheralMessage struct {
	PeripheralID string
	Payload      []byte
}

// PeripheralError is an error that occured when interacting with a Peripheral.
type PeripheralError struct {
	PeripheralID    string
	PeripheralError error
}

func (err *PeripheralError) Error() string {
	return fmt.Sprintf("[PeripheralError] %s", err.PeripheralError.Error())
}

// Listen to all the peripheral devices via a single channel.
func (peripherals PeripheralRuntime) Listen(ctx context.Context) (<-chan PeripheralMessage, <-chan error) {
	resultMessages := make(chan PeripheralMessage)
	resultErrors := make(chan error)

	for _, peripheral := range peripherals {
		payloads, errors := peripheral.Listen(ctx)

		go func(peripheral Peripheral) {
			for {
				select {
				case payload := <-payloads:
					resultMessages <- PeripheralMessage{
						peripheral.ID(),
						payload,
					}
				case err := <-errors:
					resultErrors <- &PeripheralError{
						peripheral.ID(),
						err,
					}
				case <-ctx.Done():
					return
				}
			}
		}(peripheral)
	}

	return resultMessages, resultErrors
}
