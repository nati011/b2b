package domain_handler

import (
	"context"
	"fmt"
	"strconv"

	"b2b.nati011.github.com/internal/core/application/event"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/util"
)

type ProductHandler struct {
	service product.Provider
	event   *event.Broker
}

func InitProduct(service product.Provider, event *event.Broker) error {
	handler := &ProductHandler{
		service: service,
		event:   event,
	}

	ch := event.Subscribe(util.EVENT_DISTRIBUTOR_DEACTIVATE)

	// Start a goroutine to listen for events asynchronously
	go func() {
		for msg := range ch {
			id, ok := msg.(string) // or whatever type you publish
			if !ok {
				fmt.Printf("Received unexpected event type: %T\n", msg)
				continue
			}
			did, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				return
			}
			if err := handler.service.Deactivate(context.Background(), did); err != nil {
				fmt.Printf("Failed to deactivate product %s: %v\n", id, err)
			} else {
				fmt.Printf("Deactivated product %s\n", id)
			}
		}
	}()
	return nil
}
