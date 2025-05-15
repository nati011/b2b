package db

import (
	"context"
	"time"

	port "b2b.nati011.github.com/internal/port/domain/distributor_approval"
)

type MockDistributorApproval struct {
	DistributorId int
	Verdict       bool
	Comment       string
	ReviewedAt    time.Time
	ReviewedBy    string
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) GetVerdict(ctx context.Context, id int) (port.GetVerdictResponse, error) {
	return port.GetVerdictResponse{}, nil
}

func (m *Mock) GetVerdict(ctx context.Context, id int) (port.GetVerdictResponse, error) {
	return port.GetVerdictResponse{}, nil
}
