package Delivery

import "context"

type Provider interface {
	Dispatch(ctx context.Context, req *DeliveryItemSpecificationRequest) error
	DeliverAndAcceptPayment(ctx context.Context) error
	VerifyDelivery(ctx context.Context) error
}

type DeliveryItemSpecificationRequest struct {
	ExternalId string `json:"external_id"`
}

type Delivery struct {
}

func NewPayOnDelivery() Provider {
	return &Delivery{}
}

func (p *Delivery) Dispatch(ctx context.Context, req *DeliveryItemSpecificationRequest) error {
	//create delivery
	//ensure external_id is unique
	return nil
}

func (p *Delivery) DeliverAndAcceptPayment(ctx context.Context) error {
	//mark as delivered and payment accepted
	return nil
}

func (p *Delivery) VerifyDelivery(ctx context.Context) error {
	//return delivery and payment accepted status
	return nil
}
