package domain_handler

import (
	"context"
	"errors"
	"fmt"
	"log"

	"b2b.nati011.github.com/internal/core/application/event"
	"b2b.nati011.github.com/internal/core/application/user"
	"b2b.nati011.github.com/internal/core/domain/retailer"
	"b2b.nati011.github.com/internal/core/util"
)

var (
	ErrUnknown = errors.New(" unknown error")
)

type AuthHandler struct {
	service     retailer.Provider
	userService user.Provider
	event       *event.Broker
}

func InitAuth(service retailer.Provider, userService user.Provider, event *event.Broker) error {
	handler := &AuthHandler{
		service:     service,
		userService: userService,
		event:       event,
	}

	ch := event.Subscribe(util.EVENT_RETAILER_SSO)

	// Start a goroutine to listen for events asynchronously
	go func() {
		for msg := range ch {
			payload, ok := msg.(util.EventRetailerSSOPayload) // or whatever type you publish
			if !ok {
				fmt.Printf("Received unexpected event type: %T\n", msg)
				continue
			}
			ctx := context.Background()
			_, err := handler.userService.GetByParam(ctx, &user.GetByParam{
				Email: payload.Email,
			})
			if err != nil {
				switch err {
				case user.ErrEmptyGetContent:
					handler.service.CreateAssisted(ctx, &retailer.CreateAssistedRequest{
						FirstName: payload.FirstName,
						LastName:  payload.LastName,
						Email:     payload.Email,
						Username:  payload.Username,
					})
					continue
				default:
					log.Printf("failed to get user")
					continue
				}
			}
		}
	}()
	return nil
}
