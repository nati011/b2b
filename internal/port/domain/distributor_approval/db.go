package distributorApproval

import (
	"context"
	"time"
)

type GetVerdictResponse struct {
	DistributorId int
	Verdict       bool
	Comment       string
	ReviewedBy    string
	ReviewedAt    time.Time
}

type GetAuditReportResponse struct {
	List []GetVerdictResponse
}

type RejectRequest struct {
	DistributorId int
	Comment       string
	ReviewedBy    string
}

type ApprovalRequest struct {
	DistributorId int
	Comment       string
	ReviewedBy    string
}

type Reader interface {
	GetApprovalStatus(ctx context.Context, distributorId int) (bool, error)
	GetReviewReport(ctx context.Context, distributorId int) (GetAuditReportResponse, error)
}

type Writer interface {
	Approve(ctx context.Context, req ApprovalRequest) error
	Reject(ctx context.Context, req RejectRequest) error
}

type DB interface {
	Reader
	Writer
}
