package payment

import (
	"context"
	"errors"

	factory "b2b.nati011.github.com/internal/adapter/secondary/application/payment/gateway"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	payment "b2b.nati011.github.com/internal/port/application/payment/gateway"
)

var (
	ErrUserNotSupplied                 = errors.New("oopsy, user is mandatory")
	ErrAmountNotSupplied               = errors.New("oopsy, amount is mandatory")
	ErrAmountLessThanZero              = errors.New("oopsy, amount must be greater than zero")
	ErrPaymentPartnerNotSupplied       = errors.New("oopsy, payment partner id is mandatory")
	ErrTransactionReferenceNotSupplied = errors.New("oopsy, transaction refrence is mandatory")
	ErrUnknown                         = errors.New("oopsy, unknown error has occured")
)

type CheckoutRequest struct {
	User_Id           int
	Order_id          int
	Amount            float64
	PaymentPartner_Id int
}

type CheckoutResponse struct {
	Checkout_url string
}

type VerifyResponse struct {
	Is_verified bool
}

type Provider interface {
	Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error)
	Verify(ctx context.Context, tx_ref string) (VerifyResponse, error)
}

type PaymentService struct {
	paymentPartner partner.Provider
	transaction    transaction.Provider
}

func NewPaymentService(partner partner.Provider, transaction transaction.Provider) Provider {
	return &PaymentService{
		paymentPartner: partner,
		transaction:    transaction,
	}
}

func (p *PaymentService) Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error) {

	err := p.validateAmount(req.Amount)
	if err != nil {
		return CheckoutResponse{}, err
	}
	paymentPartner, err := p.paymentPartner.GetPartnerSecret(ctx, req.PaymentPartner_Id)

	if err != nil {
		return CheckoutResponse{}, err
	}

	paymentGateway, err := factory.PaymentPartnerFactory(paymentPartner.Name)

	if err != nil {
		return CheckoutResponse{}, err
	}

	paymentInitiateRequest := payment.InitiateRequest{
		Amount:         req.Amount,
		TransactionRef: req.Order_id,
		PartnerUrl:     paymentPartner.Init_payment_url,
		PartnerSecret:  paymentPartner.Secret,
	}

	checkoutUrl, err := paymentGateway.Initiate(paymentInitiateRequest)
	if err != nil {
		return CheckoutResponse{}, err
	}

	return CheckoutResponse{
		Checkout_url: checkoutUrl,
	}, nil
}

func (p *PaymentService) Verify(ctx context.Context, tx_ref string) (VerifyResponse, error) {
	err := verifyTransactionRef(tx_ref)
	if err != nil {
		return VerifyResponse{}, err
	}
	return VerifyResponse{}, nil
}
