package user

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/distributor"
	test_container "b2b.nati011.github.com/internal/core/domain/distributor/test"
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

func Test_create_user_upon_distributor_create(t *testing.T) {
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
	resp, err := testContainer.DistributorService.GetAllUserDetail(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get user err: %v", err)
	}
	wantLen := 1
	if len(resp.List) != wantLen {
		t.Errorf("Expected len: %v Got: %v", len(resp.List), wantLen)
	}

	_, err = testContainer.UserService.Get(ctx, resp.List[0].Id)
	if err != nil {
		t.Fatalf("Failed to get user err: %v", err)
	}
}
