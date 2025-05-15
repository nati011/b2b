package distributorApproval

import (
	"context"

	port "b2b.nati011.github.com/internal/port/domain/distributor_approval"
)

type ApproveRequest struct {
	DistributorId int
}

type RejectionRequest struct {
	DistributorId int
	Comment       string
}

type Provider interface {
	Approve(ctx context.Context, id int) error
	Reject(ctx context.Context, id int) error
}

type DistributorApprovalService struct {
	DB port.DB
}

func NewDistributorApprovalService(db port.DB) Provider {
	return &DistributorApprovalService{
		DB: db,
	}
}

func (d *DistributorApprovalService) Approve(ctx context.Context, id int) error {
	return nil
}

func (d *DistributorApprovalService) Reject(ctx context.Context, id int) error {
	return nil
}
