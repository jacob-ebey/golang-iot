package azuregateway

import (
	"context"

	"github.com/amenzhinsky/iothub/iotdevice"
)

func NewAzureWriter(ctx context.Context, client *iotdevice.Client, opts ...iotdevice.SendOption) (chan []byte, chan error) {
	messages := make(chan []byte)
	errors := make(chan error)

	go func() {
		for {
			select {
			case message := <-messages:
				if err := client.SendEvent(context.Background(), message, opts...); err != nil {
					errors <- err
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return messages, errors
}

func NewAzureReader(ctx context.Context, client *iotdevice.Client) (chan []byte, chan error) {
	messages := make(chan []byte)
	errors := make(chan error)

	go func() {
		subscription, err := client.SubscribeEvents(ctx)
		if err != nil {
			errors <- err
			return
		}

		azureMessages := subscription.C()
		for {
			received := <-azureMessages
			if received != nil {
				messages <- received.Payload
			}
		}
	}()

	return messages, nil
}
