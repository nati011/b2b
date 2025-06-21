package db

import (
	"context"
	"time"

	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/distributor_approval"
)

type MockDistributorApproval struct {
	DistributorId int
	Verdict       string
	Comment       string
	ReviewedAt    time.Time
	ReviewedBy    string
}

type Mock struct {
	List []MockDistributorApproval
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) GetApprovalStatus(ctx context.Context, distributorId int) (string, error) {
	for _, i := range m.List {
		if i.DistributorId == distributorId {
			return i.Verdict, nil
		}
	}
	return "", port_commons.ErrSysNoRows
}

func (m *Mock) GetReviewReport(ctx context.Context, distributorId int) (port.GetAuditReportResponse, error) {
	var resp port.GetAuditReportResponse
	for _, i := range m.List {
		if i.DistributorId == distributorId {
			resp.List = append(resp.List, port.GetVerdictResponse{
				DistributorId: i.DistributorId,
				Verdict:       i.Verdict,
				Comment:       i.Comment,
				ReviewedAt:    i.ReviewedAt,
				ReviewedBy:    i.ReviewedBy,
			})
		}
	}
	if len(resp.List) == 0 {
		return port.GetAuditReportResponse{}, port_commons.ErrSysNoRows
	}
	return resp, nil
}

func (m *Mock) Approve(ctx context.Context, req port.ApprovalRequest) error {
	m.List = append(m.List, MockDistributorApproval{
		DistributorId: req.DistributorId,
		Verdict:       "APPROVED",
		Comment:       req.Comment,
		ReviewedAt:    time.Now(),
		ReviewedBy:    req.ReviewedBy,
	})
	return nil
}

func (m *Mock) Reject(ctx context.Context, req port.RejectRequest) error {
	m.List = append(m.List, MockDistributorApproval{
		DistributorId: req.DistributorId,
		Verdict:       "REJECTED",
		Comment:       req.Comment,
		ReviewedAt:    time.Now(),
		ReviewedBy:    req.ReviewedBy,
	})
	return nil
}
