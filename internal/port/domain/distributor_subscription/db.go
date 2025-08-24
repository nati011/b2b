package distributorApproval

import (
	"context"
)

type SubscriptonRequest struct {
	DistributorId int
}

type Reader interface {
}

type Writer interface {
	Subscribe(ctx context.Context, req SubscriptonRequest) error
}

type DB interface {
	Reader
	Writer
}
