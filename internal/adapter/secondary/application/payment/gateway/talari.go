package gateway

import port "b2b.nati011.github.com/internal/port/application/payment/gateway"

type Talari struct {
}

func NewTalari() port.Provider {
	return &Talari{}
}

func (t Talari) Initiate(req port.InitiateRequest) (string, error) {
	return "", nil
}

func (t Talari) Verify(req port.VerificationRequest) (bool, error) {
	return false, nil
}
