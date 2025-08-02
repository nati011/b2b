package onboardingapproval

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/distributor"
	test_container "b2b.nati011.github.com/internal/core/domain/distributor/test"
	distributorApproval "b2b.nati011.github.com/internal/core/domain/distributor_approval"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var testContainer test_container.TestContainer
var db *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	testContainer = test_container.NewDBIntegrationTestContainer(
		db,
	)
}

func teardown() {
	testContainer.Teardown(db)
	db_test_container.Teardown(db)
}

func Test_Init_Approval_Process_upon_registration_happypath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	in := distributor.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "dist_test",
		FirstName:   "test",
		LastName:    "test",
		Email:       "test@gmail.com",
	}
	id, err := testContainer.DistributorService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	//assert
	distApprovalStatus, err := testContainer.DistributorApprovalService.GetApprovalStatus(
		ctx, id)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	wantDistributorApprovalStatus := distributorApproval.PENDING_STATUS
	if distApprovalStatus.Status != wantDistributorApprovalStatus {
		t.Fatalf("expected err: %v, want err: %v", wantDistributorApprovalStatus, err)
	}
}
