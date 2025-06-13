package distributorApproval

import (
	"context"
	"errors"

	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/distributor_approval"
)

var (
	ErrAlreadyApproved  = errors.New("¯\\_(•_•)_/¯, already approved")
	ErrAlreadyRejected  = errors.New("¯\\_(•_•)_/¯, already rejected")
	ErrCommentMandatory = errors.New("¯\\_(•_•)_/¯, comment mandatory")
	ErrUnknown          = errors.New("¯\\_(•_•)_/¯, unknown error")
	ErrNoReview         = errors.New("¯\\_(•_•)_/¯, no reviews")
)

type ApprovalRequest struct {
	DistributorId int
}

type RejectionRequest struct {
	DistributorId int
	Comment       string
}

type Provider interface {
	Approve(ctx context.Context, req *ApprovalRequest) error
	Reject(ctx context.Context, req *RejectionRequest) error
	GetApprovalStatus(ctx context.Context, distributorId int) (bool, error)
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
	isApproved, err := d.DB.GetApprovalStatus(ctx, req.DistributorId)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return ErrUnknown
		}
	}
	if isApproved {
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

	isApproved, err := d.DB.GetApprovalStatus(ctx, req.DistributorId)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return ErrUnknown
		}
	} else {
		if !isApproved && err != port_commons.ErrSysNoRows {
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

func (d *DistributorApprovalService) GetApprovalStatus(ctx context.Context, distributorId int) (bool, error) {
	resp, err := d.DB.GetApprovalStatus(ctx, distributorId)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return false, ErrNoReview
		default:
			return false, ErrUnknown
		}
	}
	if resp {
		return true, nil
	}
	return false, nil
}
