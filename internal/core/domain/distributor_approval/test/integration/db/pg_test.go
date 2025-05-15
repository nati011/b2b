package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	distributorApproval "b2b.nati011.github.com/internal/core/domain/distributor_approval"
	"b2b.nati011.github.com/internal/core/domain/distributor_approval/test"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var testContainer test.TestContainer
var db *sql.DB
var distributorId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	testContainer = test.NewIntegrationTestContainer(
		db,
	)
	distributorId = 1
}

func teardown() {
	testContainer.Teardown(db)
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {
}

func Test_Read(t *testing.T) {
	t.Run("getApprovalStatus", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		err := testContainer.DistributorApprovalService.Reject(ctx, &distributorApproval.RejectionRequest{
			DistributorId: distributorId,
			Comment:       "test",
		})
		if err != nil {
			t.Fatalf("failed to approve distributor: %v", err)
		}
		isApproved, err := testContainer.DistributorApprovalService.GetApprovalStatus(ctx, distributorId)
		if err != nil {
			t.Fatalf("failed to get approval status err: %v", err)
		}
		wantStatus := false
		if isApproved != wantStatus {
			t.Errorf("Expected status: %v Got: %v", wantStatus, isApproved)
		}
	})
}

func Test_Write(t *testing.T) {
	t.Run("approve", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		err := testContainer.DistributorApprovalService.Reject(ctx, &distributorApproval.RejectionRequest{
			DistributorId: distributorId,
			Comment:       "test",
		})
		if err != nil {
			t.Fatalf("failed to approve distributor: %v", err)
		}
		isApproved, err := testContainer.DistributorApprovalService.GetApprovalStatus(ctx, distributorId)
		if err != nil {
			t.Fatalf("failed to get approval status err: %v", err)
		}
		wantStatus := true
		if isApproved != wantStatus {
			t.Errorf("Expected status: %v Got: %v", wantStatus, isApproved)
		}
	})

	t.Run("reject", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		err := testContainer.DistributorApprovalService.Reject(ctx, &distributorApproval.RejectionRequest{
			DistributorId: distributorId,
			Comment:       "test",
		})
		if err != nil {
			t.Fatalf("failed to approve distributor: %v", err)
		}
		isApproved, err := testContainer.DistributorApprovalService.GetApprovalStatus(ctx, distributorId)
		if err != nil {
			t.Fatalf("failed to get approval status err: %v", err)
		}
		wantStatus := false
		if isApproved != wantStatus {
			t.Errorf("Expected status: %v Got: %v", wantStatus, isApproved)
		}
	})
}
