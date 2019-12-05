package azuregateway

import (
	"context"

	"github.com/amenzhinsky/iothub/iotdevice"
	"github.com/jacob-ebey/golang-iot/core"
)

func NewAzureWriter(ctx context.Context, logger core.Logger, client *iotdevice.Client, opts ...iotdevice.SendOption) chan core.Message {
	serializer := &core.JsonSerializer{}

	writer := make(chan core.Message)
	go func() {
		for {
			message := <-writer

			payload, err := serializer.SerialiseMessage(message)
			if err != nil {
				logger.LogError(err)
				continue
			}

			if err := client.SendEvent(context.Background(), payload, opts...); err != nil {
				logger.LogError(err)
			}
		}
	}()

	return writer
}

func NewAzureReader(ctx context.Context, logger core.Logger, client *iotdevice.Client) (chan core.Message, error) {
	serializer := &core.JsonSerializer{}

	subscription, err := client.SubscribeEvents(ctx)
	if err != nil {
		return nil, err
	}

	messages := subscription.C()
	reader := make(chan core.Message)
	go func() {
		for {
			received := <-messages
			if received != nil {
				message, err := serializer.DeserialiseMessage(received.Payload)
				if err != nil {
					logger.LogError(err)
					continue
				}

				reader <- message
			}
		}
	}()

	return reader, nil
}
