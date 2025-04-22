package gateway

import port "b2b.nati011.github.com/internal/port/application/payment/gateway"

type Talari struct {
}

func NewTalari() port.Provider {
	return &Talari{}
}

func (t Talari) Initiate(req port.InitiateRequest) (port.InitatePaymentResponse, error) {
	return port.InitatePaymentResponse{}, nil
}

func (t Talari) Verify(request port.VerifyRequest) (port.VerifyResponse, error) {
	return port.VerifyResponse{}, nil
}
