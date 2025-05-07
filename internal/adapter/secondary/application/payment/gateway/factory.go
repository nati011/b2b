package gateway

import (
	"errors"

	port "b2b.nati011.github.com/internal/port/application/payment/gateway"
)

var (
	ErrSysInstaceNotFound = errors.New("instance not found")
)

func PaymentPartnerFactory(partnerName string) (port.Provider, error) {
	switch partnerName {
	case "talari":
		return Talari{}, nil
	case "chapa":
		return Chapa{}, nil
	case "payOnDelivery":
		return PayOnDelivery{}, nil
	default:
		return nil, ErrSysInstaceNotFound
	}
}
