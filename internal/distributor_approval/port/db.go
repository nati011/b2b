package distributorApproval

import (
	"context"
	"time"
)

type ApprovalResult struct {
	Status string
}

type GetVerdictResponse struct {
	DistributorId int
	Verdict       string
	Comment       string
	ReviewedBy    string
	ReviewedAt    time.Time
}

type GetAuditReportResponse struct {
	List []GetVerdictResponse
}

type ReviewChangeRequest struct {
	DistributorId int
	Verdict       string
	Comment       string
	ReviewedBy    string
}

type InitRequest struct {
	DistributorId int
	Status        string
}

type Reader interface {
	GetApprovalStatus(ctx context.Context, distributorId int) (string, error)
	GetReviewReport(ctx context.Context, distributorId int) (GetAuditReportResponse, error)
}

type Writer interface {
	ChangeDistributorReview(ctx context.Context, req ReviewChangeRequest) error
	Init(ctx context.Context, req InitRequest) error
}

type DB interface {
	Reader
	Writer
}
