package user

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	user "b2b.nati011.github.com/internal/core/application/user"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var testContainer user.TestContainer
var db *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	testContainer = user.NewIntegrationTestContainer(
		db,
	)
}

func teardown() {
	testContainer.TeardownIntegrationTestContainer()
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {}

func Test_write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName:  "natnael jemaneh asefa",
			LastName:   "natnael",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			ExternalId: "123",
		}
		id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		_, err = testContainer.UserService.Get(ctx, id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}
	})

	t.Run("remove", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName:  "natnael jemaneh asefa",
			LastName:   "natnael",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			ExternalId: "123",
		}
		id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//remove
		err = testContainer.UserService.Remove(ctx, id)
		if err != nil {
			t.Fatalf("Failed to remove user err: %v", err)
		}
		_, err = testContainer.UserService.Get(ctx, id)
		wantErr := user.ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v got: %v", wantErr, err)
		}
	})
}

func Test_read(t *testing.T) {
	t.Run("get_all", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName: "natnael jemaneh asefa",
			LastName:  "natnael",
			Email:     "natnaeljemaneh001@gmail.com",
			Phone:     "+251949184879",
			Username:  "test",
			DOB:       parsedTime,

			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		in_2 := user.CreateRequest{
			FirstName: "natnael jemaneh asefa",
			LastName:  "natnael",
			Email:     "natnaeljemaneh002@gmail.com",
			Phone:     "+251949184889",
			Username:  "test_009",
			DOB:       parsedTime,

			ExternalId: "123",
		}
		_, err = testContainer.UserService.Create(ctx, &in_2)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		pagination := user.Pagination{
			Limit:  1,
			Offset: 0,
		}
		//get
		resp, err := testContainer.UserService.GetAll(ctx, &pagination)
		if err != nil {
			t.Fatalf("Failed to getAll err: %v", err)
		}
		wantLen := pagination.Limit
		if wantLen != len(resp.List) && wantLen < len(resp.List) {
			t.Errorf("Expected length: %v Want: %v", wantLen, len(resp.List))
		}
	})

	t.Run("get_by_param", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName:  "natnael jemaneh asefa",
			LastName:   "natnael",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		inParam := &user.GetByParam{
			Username: "test",
		}
		pagination := user.Pagination{
			Limit:  1,
			Offset: 0,
		}
		response, err := testContainer.UserService.GetByParam(ctx, inParam, &pagination)
		if err != nil {
			t.Fatalf("Failed to get user by param err: %v", err)
		}
		expecetdLen := 1
		if len(response.List) != expecetdLen {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, len(response.List))
		}
	})

	t.Run("update", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02", "2024-09-20")
		in := user.CreateRequest{
			FirstName:  "natnael jemaneh asefa",
			LastName:   "natnael",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		in_updateParsedTime, _ := time.Parse("2006-01-02", "2024-09-19")
		in_update := &user.UpdateRequest{
			Id:         user_id,
			FirstName:  "test",
			LastName:   "natnael",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        in_updateParsedTime,
			ExternalId: "123",
		}
		got, err := testContainer.UserService.Update(ctx, in_update)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}
		if got.FirstName != in_update.FirstName {
			t.Errorf("Expected : %v Got: %v", in_update.DOB, got.DOB)
		}
		if got.Email != in_update.Email {
			t.Errorf("Expected : %v Got: %v", in_update.Email, got.Email)
		}
		if got.Phone != in_update.Phone {
			t.Errorf("Expected : %v Got: %v", in_update.Phone, got.Phone)
		}
		if got.Username != in_update.Username {
			t.Errorf("Expected : %v Got: %v", in_update.Username, got.Username)
		}
		if got.DOB != in_update.DOB {
			t.Errorf("Expected : %v Got: %v", in_update.DOB, got.DOB)
		}
		if got.ExternalId != in_update.ExternalId {
			t.Errorf("Expected : %v Got: %v", in_update.ExternalId, got.ExternalId)
		}
	})

}
