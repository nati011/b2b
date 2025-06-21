package db

import (
	"context"
	"database/sql"

	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/domain/distributor_approval"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(DB *sql.DB) port.DB {
	return &Postgres{
		Pool: DB,
	}
}

func (m *Postgres) GetApprovalStatus(ctx context.Context, distributorId int) (string, error) {
	var approval_result port.ApprovalResult
	query := "SELECT * FROM public.get_approval_status_by_distributor_id($1);"
	result := []any{&approval_result.Status}

	args := []any{distributorId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(m.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return "", err
	}
	approval_result.Status = *result[0].(*string)
	return approval_result.Status, nil
}

func (m *Postgres) GetReviewReport(ctx context.Context, distributorId int) (port.GetAuditReportResponse, error) {
	return port.GetAuditReportResponse{}, nil
}

func (m *Postgres) Approve(ctx context.Context, req port.ApprovalRequest) error {
	query := "SELECT * FROM public.approve_distributor_review($1, $2, $3);"

	args := []any{
		req.DistributorId,
		req.Comment,
		req.ReviewedBy}

	if err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(m.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery(); err != nil {
		return err
	}
	return nil
}

func (m *Postgres) Reject(ctx context.Context, req port.RejectRequest) error {
	query := "SELECT * FROM public.reject_distributor_review($1, $2, $3);"

	args := []any{
		req.DistributorId,
		req.Comment,
		req.ReviewedBy}

	if err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(m.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery(); err != nil {
		return err
	}
	return nil
}
