package payment

import (
	"context"
	"errors"
<<<<<<< HEAD
	"log"

	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/application/user"
	payment_processor "b2b.nati011.github.com/internal/core/domain/paymentProcessor"
=======

	factory "b2b.nati011.github.com/internal/adapter/secondary/application/payment/gateway"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	payment "b2b.nati011.github.com/internal/port/application/payment/gateway"
>>>>>>> c481e966 (init handle multiple payment gateway)
)

var (
	ErrUserNotFound                    = errors.New("oopsy, user is found")
	ErrAmountNotSupplied               = errors.New("oopsy, amount is mandatory")
	ErrAmountLessThanZero              = errors.New("oopsy, amount must be greater than zero")
	ErrPaymentPartnerNotSupplied       = errors.New("oopsy, payment partner id is mandatory")
	ErrTransactionReferenceNotSupplied = errors.New("oopsy, transaction refrence is mandatory")
	ErrUnknown                         = errors.New("oopsy, unknown error has occured")
)

type CheckoutRequest struct {
<<<<<<< HEAD
	Amount           int64
	PaymentPartnerId int
=======
	User_Id           int
	Order_id          int
	Amount            float64
	PaymentPartner_Id int
>>>>>>> c481e966 (init handle multiple payment gateway)
}

type CheckoutResponse struct {
	Checkout_url string
}

type Provider interface {
	Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error)
	Verify(ctx context.Context, gateway_id int, tx_ref string) (bool, error)
	Callback(ctx context.Context, gateway_id int, tx_ref string)
}

type PaymentService struct {
<<<<<<< HEAD
	UserService        user.Provider
	PartnerService     partner.Provider
	TransactionService transaction.Provider
	PaymentProcessor   payment_processor.Provider
}

func NewPaymentService(
	ps partner.Provider,
	ts transaction.Provider) Provider {
	return &PaymentService{
		PartnerService:     ps,
		TransactionService: ts,
=======
	paymentPartner partner.Provider
	transaction    transaction.Provider
}

func NewPaymentService(partner partner.Provider, transaction transaction.Provider) Provider {
	return &PaymentService{
		paymentPartner: partner,
		transaction:    transaction,
>>>>>>> c481e966 (init handle multiple payment gateway)
	}
}

func (p *PaymentService) Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error) {
<<<<<<< HEAD
=======

>>>>>>> c481e966 (init handle multiple payment gateway)
	err := p.validateAmount(req.Amount)
	if err != nil {
		return CheckoutResponse{}, err
	}
<<<<<<< HEAD
	return CheckoutResponse{}, nil
=======
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
>>>>>>> c481e966 (init handle multiple payment gateway)
}

func (p *PaymentService) Verify(ctx context.Context, gateway_id int, tx_ref string) (bool, error) {
	return false, nil
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
		p.PaymentProcessor.Process(ctx, tx_ref)
	}
}
