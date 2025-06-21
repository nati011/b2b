package distributorApproval

import (
	"context"
	"errors"

	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/distributor_approval"
)

var (
	ErrAlreadyApproved  = errors.New(" already approved")
	ErrAlreadyRejected  = errors.New(" already rejected")
	ErrCommentMandatory = errors.New(" comment mandatory")
	ErrUnknown          = errors.New(" unknown error")
	ErrNoReview         = errors.New(" no reviews")
)

var (
	REJECTED_STATUS = "REJECTED"
	PENDING_STATUS  = "PENDING"
	APPROVED_STATUS = "APPROVED"
)

type ApprovalRequest struct {
	DistributorId int
}

type RejectionRequest struct {
	DistributorId int
	Comment       string
}

type ApprovalStatusResponse struct {
	Status string
}

type Provider interface {
	Approve(ctx context.Context, req *ApprovalRequest) error
	Reject(ctx context.Context, req *RejectionRequest) error
	GetApprovalStatus(ctx context.Context, distributorId int) (ApprovalStatusResponse, error)
}

type DistributorApprovalService struct {
	DB port.DB
}

func NewDistributorApprovalService(db port.DB) Provider {
	return &DistributorApprovalService{
		DB: db,
	}
}

func (d *DistributorApprovalService) Approve(ctx context.Context, req *ApprovalRequest) error {
	status, err := d.DB.GetApprovalStatus(ctx, req.DistributorId)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return ErrUnknown
		}
	}
	if status == APPROVED_STATUS {
		return ErrAlreadyApproved
	}
	err = d.DB.Approve(ctx, port.ApprovalRequest{
		DistributorId: req.DistributorId,
	})
	if err != nil {
		return ErrUnknown
	}
	return nil
}

func (d *DistributorApprovalService) Reject(ctx context.Context, req *RejectionRequest) error {
	if req.Comment == "" {
		return ErrCommentMandatory
	}

	status, err := d.DB.GetApprovalStatus(ctx, req.DistributorId)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return ErrUnknown
		}
	} else {
		if status == REJECTED_STATUS && err != port_commons.ErrSysNoRows {
			return ErrAlreadyRejected
		}
	}

	err = d.DB.Reject(ctx, port.RejectRequest{
		DistributorId: req.DistributorId,
		Comment:       req.Comment,
	})
	if err != nil {
		return ErrUnknown
	}
	return nil
}

func (d *DistributorApprovalService) GetApprovalStatus(ctx context.Context, distributorId int) (ApprovalStatusResponse, error) {
	resp, err := d.DB.GetApprovalStatus(ctx, distributorId)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return ApprovalStatusResponse{}, ErrNoReview
		default:
			return ApprovalStatusResponse{}, ErrUnknown
		}
	}
	return ApprovalStatusResponse{Status: resp}, nil
}
