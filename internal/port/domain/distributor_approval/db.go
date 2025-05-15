package distributorApproval

import (
	"context"
)

type GetVerdictResponse struct {
	DistributorId int
	Verdict       bool
	Comment       string
}

type GetAuditReportResponse struct {
	List []GetVerdictResponse
}

type CreateResponse struct {
	DistributorId int
	Verdict       bool
	Comment       string
	ReviewedBy    string
}

type Reader interface {
	GetVerdict(ctx context.Context, id int) (GetVerdictResponse, error)
	GetAuditReport(ctx context.Context) (GetAllVerdictResponse, error)
}

type Writer interface {
	Approve(ctx context.Context, id int) error
	Reject(ctx context.Context, id int) error
}

type DB interface {
	Reader
	Writer
}
