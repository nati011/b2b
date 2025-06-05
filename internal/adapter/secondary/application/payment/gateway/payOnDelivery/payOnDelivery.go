package gateway

import port "b2b.nati011.github.com/internal/port/application/payment/gateway"

type PayOnDelivery struct {
}

func NewPayOnDelivery() port.Provider {
	return &PayOnDelivery{}
}

func (p PayOnDelivery) Initiate(request port.InitiateRequest) (string, error) {
	return "", nil
}

func (p PayOnDelivery) Verify(request port.VerificationRequest) (bool, error) {
	return false, nil
}
