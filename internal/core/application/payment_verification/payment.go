package payment_verification

import (
	"context"
	"errors"
	"log"
	"time"

	factory "b2b.nati011.github.com/internal/adapter/secondary/application/payment/gateway"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/domain/order"
	port "b2b.nati011.github.com/internal/port/application/payment/db"
	payment "b2b.nati011.github.com/internal/port/application/payment/gateway"
)

var (
	ErrUserNotFound                    = errors.New("oopsy, user is found")
	ErrAmountNotSupplied               = errors.New("oopsy, amount is mandatory")
	ErrAmountLessThanZero              = errors.New("oopsy, amount must be greater than zero")
	ErrPaymentPartnerNotSupported      = errors.New("oopsy, payment partner id is mandatory")
	ErrTransactionReferenceNotSupplied = errors.New("oopsy, transaction refrence is mandatory")
	ErrUnknown                         = errors.New("oopsy, unknown error has occured")
	ErrCannotProceedWithPaymentPartner = errors.New("oopsy, cannot proceed with payment partner")
)

type CheckoutRequest struct {
	OrderId          int
	PaymentPartnerId int
}

type GetPaymentResponse struct {
	Id             int
	Date           time.Time
	Amount         float64
	OrderId        int
	PartnerId      int
	TransactionRef string
}

type CheckoutResponse struct {
	Checkout_url string
}

type Provider interface {
	Verify(ctx context.Context, PaymentPartnerId int, tx_ref string) (bool, error)
	Callback(ctx context.Context, PaymentPartnerId int, tx_ref string)
}

type PaymentService struct {
	db             port.DB
	paymentPartner partner.Provider
	transaction    transaction.Provider
	order          order.Provider
}

type CreatePaymentRequest struct {
	PaymentPartnerId int
	OrderId          int
}

func NewPaymentVerificationService(DB port.DB, partner partner.Provider, transaction transaction.Provider, order order.Provider) Provider {
	return &PaymentService{
		db:             DB,
		paymentPartner: partner,
		transaction:    transaction,
		order:          order,
	}
}

func (p *PaymentService) getPayment(ctx context.Context, tx_ref string) (GetPaymentResponse, error) {
	payment, err := p.db.GetByTransactionRef(ctx, tx_ref)
	if err != nil {
		return GetPaymentResponse{}, err
	}

	return (GetPaymentResponse)(payment), err

}

func (p *PaymentService) Verify(ctx context.Context, gateway_id int, tx_ref string) (bool, error) {
	paymentPartner, err := p.paymentPartner.GetPartnerSecret(ctx, gateway_id)
	if err != nil {
		switch err {
		case partner.ErrIdNotFound:
			return false, ErrPaymentPartnerNotSupported
		default:
			return false, ErrUnknown
		}
	}

	paymentGateway, err := factory.PaymentPartnerFactory(paymentPartner.Name)
	if err != nil {
		return false, ErrCannotProceedWithPaymentPartner
	}

	paymentVerificationRequest := payment.VerificationRequest{
		PartnerUrl:     paymentPartner.BaseURL,
		TransactionRef: tx_ref,
		PartnerSecret:  paymentPartner.Secret,
	}

	is_verified, err := paymentGateway.Verify(paymentVerificationRequest)
	if err != nil {
		return false, ErrUnknown
	}
	return is_verified, nil
}

func (p *PaymentService) Callback(ctx context.Context, gateway_id int, tx_ref string) {
	is_verified, err := p.Verify(ctx, gateway_id, tx_ref)
	if err != nil {
		switch err {
		default:
			log.Printf("failed to process incoming callback tx_ref: %v, gateway_id: %v", tx_ref, gateway_id)
		}
	}

	if is_verified {
		payment, err := p.getPayment(ctx, tx_ref)
		if err != nil {
			log.Printf("Error Occured while fetching payment: %v", err)
		}

		err = p.transaction.UpdateByTransactionRef(ctx, &transaction.UpdateByTransactionRefRequest{
			TransactionRef: tx_ref,
			Status:         transaction.COMPLETED_STATUS,
		})
		if err != nil {
			log.Printf("Error Occured while updating transaction status: %v", err)
		}

		_, err = p.order.UpdatePaymentStatus(ctx, &order.UpdateRequest{
			Id:            payment.OrderId,
			PaymentStatus: order.PAYMENT_ACCEPTED_STATUS,
		})
		if err != nil {
			log.Printf("Error Occured while updating payment status: %v", err)
		}
	}
}
