package distributor_subscription

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/distributor_subscription"
)

var (
	ErrEmptyGetContent                         = errors.New(" empty conent")
	ErrNameMandatory                           = errors.New(" name mandatory")
	ErrDescriptionMandatory                    = errors.New(" desc mandatory")
	ErrTermCannotBeZero                        = errors.New(" term cannot be zero")
	ErrTermCannotBeNegative                    = errors.New(" term cannot be negative")
	ErrPriceCannotBeZero                       = errors.New(" price cannot be zero")
	ErrPriceCannotBeNegative                   = errors.New(" price cannot be negative")
	ErrSubscriptionAlreadyExistsForDistributor = errors.New(" subscription already exists for distributor")
	ErrSubscriptionNotFoundForDistributor      = errors.New(" subscription not found for distributor")
	ErrSubscriptionNotFound                    = errors.New(" subscription not found")
	ErrPaymentPartnerIdNotSupported            = errors.New(" payment partner not supported")
	ErrSubsctiptionAlreadyActive               = errors.New(" subscription already active")
	ErrUnknown                                 = errors.New(" unknown error")
)

var (
	STATUS_SUBSCRIPTION_ACTIVE  = "active"
	STATUS_SUBSCRIPTION_EXPIRED = "expired"
)

type PlaceRequest struct {
	SubscriptionPlanId int
	DistributorId      int
	PaymentPartnerId   int
}

type GetSubscriptionResponse struct {
	Id                 int       `json:"id"`
	SubscriptionPlanId int       `json:"subscription_plan_id"`
	DistributorId      int       `json:"distributor_id"`
	Status             string    `json:"status"`
	CreatedDate        time.Time `json:"created_date"`
}

type GetAllSubscriptionResponse struct {
	List []GetSubscriptionResponse
}

type SubscribeResponse struct {
	Id          int    `json:"id"`
	CheckoutUrl string `json:"checkout_url"`
	TxRef       string `json:"tx_ref"`
}

type GetSubscriptionPlanResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	TermInMonth int     `json:"term_in_month"`
	Description string  `json:"desc"`
}

type GetAllSubscriptionPlanResponse struct {
	List []GetSubscriptionPlanResponse `json:"list"`
}

type CreatePlanRequest struct {
	Name        string
	Price       float64
	TermInMonth int
	Description string
}

// TODO: renew subscription api
// TODO: expired subscription notification
// TODO: acitvate and deactivate plans
// TODO: make name, price, termInMonth unique
// Note: instead of updating plans, deactivate plan of choice and create a new one instead
type Prodvider interface {
	CreatePlan(ctx context.Context, req *CreatePlanRequest) (int, error)
	GetPlan(ctx context.Context, id int) (GetSubscriptionPlanResponse, error)
	GetAllPlan(ctx context.Context) (GetAllSubscriptionPlanResponse, error)
	Place(ctx context.Context, req *PlaceRequest) (SubscribeResponse, error)
	InitPayment(ctx context.Context, subscriptionId int) (SubscribeResponse, error)
	GetSubscription(ctx context.Context, subscriptionId int) (GetSubscriptionResponse, error)
	GetSubscriptions(ctx context.Context) (GetAllSubscriptionResponse, error)
	GetSubscriptionByDistributorId(ctx context.Context, distId int) (GetSubscriptionResponse, error)
	ExpireSubscription(ctx context.Context, subscriptionId int) error
	RenewSubscription(ctx context.Context, subscriptionId int) (SubscribeResponse, error)
}

type DistributorSubscriptionService struct {
	DB              port.DB
	Checkout        checkout.Provider
	Payment_partner payment_partner.Provider
}

func NewDistributorSubscriptionService(
	db port.DB,
	checkout checkout.Provider,
	partner payment_partner.Provider) Prodvider {
	return &DistributorSubscriptionService{
		DB:              db,
		Checkout:        checkout,
		Payment_partner: partner,
	}
}

func (d *DistributorSubscriptionService) CreatePlan(ctx context.Context, req *CreatePlanRequest) (int, error) {
	if err := d.validateCreate(req); err != nil {
		return 0, err
	}
	id, err := d.DB.CreatePlan(ctx, &port.CreatePlanRequest{
		Name:        req.Name,
		Price:       req.Price,
		TermInMonth: req.TermInMonth,
		Description: req.Description,
	})
	if err != nil {
		return 0, ErrUnknown
	}
	return id, nil
}

func (d *DistributorSubscriptionService) GetAllPlan(ctx context.Context) (GetAllSubscriptionPlanResponse, error) {
	resp, err := d.DB.GetAllPlan(ctx)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllSubscriptionPlanResponse{}, ErrEmptyGetContent
		default:
			log.Printf("failed to get distributor subscription plans err: %v", err.Error())
			return GetAllSubscriptionPlanResponse{}, ErrUnknown
		}
	}
	response := GetAllSubscriptionPlanResponse{}
	for _, i := range resp.List {
		response.List = append(response.List, GetSubscriptionPlanResponse{
			Id:          i.Id,
			Name:        i.Name,
			Price:       i.Price,
			TermInMonth: i.TermInMonth,
			Description: i.Description,
		})
	}
	return response, nil
}

func (d *DistributorSubscriptionService) GetPlan(ctx context.Context, id int) (GetSubscriptionPlanResponse, error) {
	resp, err := d.DB.GetPlan(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetSubscriptionPlanResponse{}, ErrEmptyGetContent
		default:
			log.Printf("failed to get distributor subscription plans err: %v", err.Error())
			return GetSubscriptionPlanResponse{}, ErrUnknown
		}
	}
	return GetSubscriptionPlanResponse(resp), nil
}

func (d *DistributorSubscriptionService) Place(ctx context.Context, req *PlaceRequest) (SubscribeResponse, error) {
	//validate subscription plan exists
	plan_resp, err := d.GetPlan(ctx, req.SubscriptionPlanId)
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return SubscribeResponse{}, ErrEmptyGetContent
		default:
			log.Printf("failed to get plan: %v", err)
			return SubscribeResponse{}, ErrUnknown
		}
	}

	//validate paymentPartnerId
	_, err = d.Payment_partner.Get(ctx, req.PaymentPartnerId)
	if err != nil {
		switch err {
		case payment_partner.ErrIdNotFound:
			return SubscribeResponse{}, ErrPaymentPartnerIdNotSupported
		default:
			log.Printf("failed to get payment partner err: %v", err)
			return SubscribeResponse{}, ErrUnknown
		}
	}

	//validate if subscription already exists
	sub_resp, err := d.GetSubscriptionByDistributorId(ctx, req.DistributorId)
	wantErr := ErrEmptyGetContent
	if err != wantErr {
		switch err {
		case nil:
			return SubscribeResponse{}, ErrSubscriptionAlreadyExistsForDistributor
		default:
			log.Printf("failed to get subescription by distributorId err: %v", err)
			return SubscribeResponse{}, ErrUnknown
		}
	}

	//create subscription record
	sub_id, err := d.DB.Place(ctx, &port.PlaceRequest{
		PaymentPartnerId:   req.PaymentPartnerId,
		SubscriptionPlanId: sub_resp.Id,
		DistributorId:      req.DistributorId,
		Status:             STATUS_SUBSCRIPTION_ACTIVE,
	})
	if err != nil {
		switch err {
		default:
			log.Printf("failed to place err: %v", err)
			return SubscribeResponse{}, ErrUnknown
		}
	}

	//init subscription checkout and return
	pay_resp, err := d.Checkout.Checkout(ctx, &checkout.CheckoutRequest{
		OrderId:          sub_id,
		Amount:           plan_resp.Price,
		PaymentPartnerId: req.PaymentPartnerId,
	})
	if err != nil {
		switch err {
		default:
			log.Printf("failed to checkout err: %v", err)
			return SubscribeResponse{}, ErrUnknown
		}
	}
	return SubscribeResponse{
		Id:          sub_id,
		CheckoutUrl: pay_resp.CheckoutUrl,
		TxRef:       pay_resp.TransactionRef,
	}, nil
}

func (d *DistributorSubscriptionService) RenewSubscription(ctx context.Context, subscriptionId int) (SubscribeResponse, error) {
	//validate subscription exists
	sub, err := d.GetSubscription(ctx, subscriptionId)
	if err != nil {
		return SubscribeResponse{}, ErrSubscriptionNotFound
	}

	//validate if subscription has expired
	if strings.Split(sub.Status, "") == nil {
		log.Printf("failed to get subscription status")
		return SubscribeResponse{}, ErrUnknown
	}

	if sub.Status == STATUS_SUBSCRIPTION_ACTIVE {
		return SubscribeResponse{}, ErrSubsctiptionAlreadyActive
	}

	resp, err := d.InitPayment(ctx, sub.DistributorId)
	if err != nil {
		log.Printf("faield to init payment %v", err.Error())
		return SubscribeResponse{}, ErrUnknown
	}
	return resp, nil
}

func (d *DistributorSubscriptionService) InitPayment(ctx context.Context, distId int) (SubscribeResponse, error) {
	//get subscription
	sub, err := d.GetSubscriptionByDistributorId(ctx, distId)
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return SubscribeResponse{}, ErrSubscriptionNotFoundForDistributor
		default:
			log.Printf("failed to place err: %v", err)
			return SubscribeResponse{}, ErrUnknown
		}
	}

	//validate subscription plan exists
	plan_resp, err := d.GetPlan(ctx, sub.SubscriptionPlanId)
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return SubscribeResponse{}, ErrEmptyGetContent
		default:
			log.Printf("failed to get plan: %v", err)
			return SubscribeResponse{}, ErrUnknown
		}
	}

	//init subscription checkout and return
	pay_resp, err := d.Checkout.Checkout(ctx, &checkout.CheckoutRequest{
		OrderId: sub.Id,
		Amount:  plan_resp.Price,
		// PaymentPartnerId: plan_resp.PaymentPartnerId,
	})
	if err != nil {
		switch err {
		default:
			log.Printf("failed to checkout err: %v", err)
			return SubscribeResponse{}, ErrUnknown
		}
	}
	return SubscribeResponse{
		Id:          sub.Id,
		CheckoutUrl: pay_resp.CheckoutUrl,
		TxRef:       pay_resp.TransactionRef,
	}, nil
}

func (d *DistributorSubscriptionService) GetSubscription(ctx context.Context, subId int) (GetSubscriptionResponse, error) {
	resp, err := d.DB.GetSubscription(ctx, subId)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetSubscriptionResponse{}, ErrEmptyGetContent
		default:
			log.Printf("failed to get sub: %v", err)
			return GetSubscriptionResponse{}, ErrUnknown
		}
	}

	return GetSubscriptionResponse{
		Id:                 resp.Id,
		SubscriptionPlanId: resp.SubscriptionPlanId,
		DistributorId:      resp.DistributorId,
		Status:             resp.Status,
	}, nil
}

func (d *DistributorSubscriptionService) GetSubscriptions(ctx context.Context) (GetAllSubscriptionResponse, error) {
	response := GetAllSubscriptionResponse{}
	resp, err := d.DB.GetAllSubscriptions(ctx)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllSubscriptionResponse{}, ErrEmptyGetContent
		default:
			log.Printf("failed to get all subs: %v", err)
			return GetAllSubscriptionResponse{}, ErrUnknown
		}
	}
	for _, i := range resp.List {
		response.List = append(response.List, GetSubscriptionResponse{
			Id:                 i.Id,
			SubscriptionPlanId: i.SubscriptionPlanId,
			DistributorId:      i.DistributorId,
			Status:             i.Status,
		})
	}
	return response, nil
}

func (d *DistributorSubscriptionService) GetSubscriptionByDistributorId(ctx context.Context, distId int) (GetSubscriptionResponse, error) {
	resp, err := d.DB.GetSubscriptionByDistributorId(ctx, distId)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetSubscriptionResponse{}, ErrEmptyGetContent
		default:
			log.Printf("failed to get subscription by distributorId  err: %v", err)
			return GetSubscriptionResponse{}, ErrUnknown
		}
	}
	return GetSubscriptionResponse{
		Id:                 resp.Id,
		SubscriptionPlanId: resp.SubscriptionPlanId,
		DistributorId:      resp.DistributorId,
		Status:             resp.Status,
	}, nil
}

func (d *DistributorSubscriptionService) ExpireSubscription(ctx context.Context, subscriptionId int) error {
	err := d.DB.UpdateStatus(ctx, &port.UpdateSubscriptionRequest{
		Id:     subscriptionId,
		Status: STATUS_SUBSCRIPTION_EXPIRED,
	})
	if err != nil {
		log.Printf("failed to expire subsctiption")
		return ErrUnknown
	}
	return nil
}
