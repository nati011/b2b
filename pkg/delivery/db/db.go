package Delivery

import "context"

type GetDeliveryResponse struct {
	UUID              string `json:"uuid"`
	ExternalId        string `json:"external_id"`
	CreatedDate       string `json:"created_date"`
	IsDelivered       bool   `json:"is_delivered"`
	IsPaymentAccepted bool   `json:"is_payment_accepted"`
}

type CreateDeliveryRequest struct {
	ExternalId string `json:"external_id"`
}

type DB interface {
	CreateDelivery(ctx context.Context, req *CreateDeliveryRequest)
	Deliver(ctx context.Context, uuid string)
	AcceptPayment(ctx context.Context, uuid string)
	GetDelivery(ctx context.Context, uuid string) (GetDeliveryResponse, error)
}
