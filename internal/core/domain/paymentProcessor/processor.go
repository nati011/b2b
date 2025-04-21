package payment_processor

import (
	"b2b.nati011.github.com/internal/core/domain/order"
)

type Provider struct {
	orderService order.Provider
}

func (p *Provider) Process(tx_ref string) {
	//change order payment status to PAYMENT_ACCEPTED
}
