package payment_verification

import (
	"context"
	"errors"
	"log"
	"time"

	factory "b2b.nati011.github.com/internal/adapter/secondary/application/payment/gateway"
	payment "b2b.nati011.github.com/internal/core/application/payment"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/domain/order"
	payment_gateway "b2b.nati011.github.com/internal/port/application/payment/gateway"
)

var (
	ErrUserNotFound                    = errors.New("oopsy, user is found")
	ErrAmountNotSupplied               = errors.New("oopsy, amount is mandatory")
	ErrAmountLessThanZero              = errors.New("oopsy, amount must be greater than zero")
	ErrPaymentPartnerNotSupported      = errors.New("oopsy, payment partner id is not supported")
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
	Verify(ctx context.Context, txRef string) (bool, error)
	Callback(ctx context.Context, PaymentPartnerId int, txRef string)
}

type PaymentVerificationService struct {
	paymentService payment.Provider
	paymentPartner partner.Provider
	transaction    transaction.Provider
	order          order.Provider
}

type CreatePaymentRequest struct {
	PaymentPartnerId int
	OrderId          int
}

func NewPaymentVerificationService(paymentService payment.Provider, partner partner.Provider, transaction transaction.Provider, order order.Provider) Provider {
	return &PaymentVerificationService{
		paymentService: paymentService,
		paymentPartner: partner,
		transaction:    transaction,
		order:          order,
	}
}

func (p *PaymentVerificationService) getPayment(ctx context.Context, txRef string) (GetPaymentResponse, error) {
	payment, err := p.paymentService.GetByTransactionRef(ctx, txRef)
	if err != nil {
		return GetPaymentResponse{}, err
	}
	return GetPaymentResponse{
		Id:             payment.Id,
		Date:           payment.Date,
		Amount:         payment.Amount,
		OrderId:        payment.OrderId,
		PartnerId:      payment.PartnerId,
		TransactionRef: payment.TransactionRef,
	}, err
}

func (p *PaymentVerificationService) Verify(ctx context.Context, txRef string) (bool, error) {
	if err := validateTxRef(txRef); err != nil {
		return false, err
	}
	payment, err := p.transaction.GetByParam(ctx, &transaction.GetByParamRequest{TxRef: txRef})
	if err != nil {
		return false, ErrUnknown
	}
	paymentPartnerId := payment.List[0].PartnerId
	paymentPartner, err := p.paymentPartner.Get(ctx, paymentPartnerId)
	if err != nil {
		switch err {
		case partner.ErrIdNotFound:
			return false, ErrPaymentPartnerNotSupported
		default:
			return false, ErrUnknown
		}
	}
	if paymentPartner.Status != partner.ACTIVE_STATUS {
		return false, ErrPaymentPartnerNotSupported
	}

	paymentPartnerSecret, err := p.paymentPartner.GetPartnerSecret(ctx, paymentPartnerId)
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

	is_verified, err := paymentGateway.Verify(payment_gateway.VerificationRequest{
		PartnerUrl:     paymentPartner.BaseURL,
		TransactionRef: txRef,
		PartnerSecret:  paymentPartnerSecret.Secret,
	})
	if err != nil {
		return false, ErrUnknown
	}
	return is_verified, nil
}

func (p *PaymentVerificationService) Callback(ctx context.Context, gatewayId int, txRef string) {
	is_verified, err := p.Verify(ctx, txRef)
	if err != nil {
		switch err {
		default:
			log.Printf("failed to process incoming callback txRef: %v, gateway_id: %v", txRef, gatewayId)
		}
	}

	if is_verified {
		payment, err := p.getPayment(ctx, txRef)
		if err != nil {
			log.Printf("Error Occured while fetching payment: %v", err)
		}

		err = p.transaction.UpdateByTransactionRef(ctx, &transaction.UpdateByTransactionRefRequest{
			TransactionRef: txRef,
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
