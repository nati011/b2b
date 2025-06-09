package transaction

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	test_container "b2b.nati011.github.com/internal/core/application/transaction/test"
	"b2b.nati011.github.com/internal/core/application/user"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var container test_container.TestContainer
var db *sql.DB
var user_id int
var PartnerId int

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	container = test_container.NewDBIntegrationTestContainer(db)
	ctx := context.Background()
	var err error
	PartnerId, err = container.PartnerService.Create(ctx, &partner.CreateRequest{
		Name:          "test",
		Icon:          "test",
		BaseURL:       "test",
		Secret:        "test",
		PaymentMethod: payment_partner.PAYMENT_METHOD_DIGITAL,
	})
	if err != nil {
		panic("failed to create err partner")
	}
	//create user
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := user.CreateRequest{
		FirstName:  "natnael jemaneh asefa",
		LastName:   "test",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		ExternalId: "123",
	}
	user_id, _ = container.UserService.Create(ctx, &in)
}

func teardown() {
	container.TeardownIntegrationTestContainer(db)
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {

}

func Test_write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &transaction.CreateRequest{
			Amount:    1,
			PartnerId: PartnerId,
			TxRef:     "test",
		}
		id, err := container.TransactionService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check
		resp, err := container.TransactionService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.Id)
		}
	})

	t.Run("update_status", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &transaction.CreateRequest{
			Amount:    1,
			PartnerId: PartnerId,
			TxRef:     "test",
		}
		id, err := container.TransactionService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		update_in := transaction.UpdateRequest{
			Id:     id,
			Status: "test_two",
		}
		err = container.TransactionService.UpdateStatus(ctx, &update_in)
		if err != nil {
			t.Errorf("Failed to update status err: %v", err)
		}
		//check
		resp, err := container.TransactionService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if resp.Status != update_in.Status {
			t.Errorf("Expected status: %v Got: %v", update_in.Status, resp.Status)
		}
	})
}

func Test_read(t *testing.T) {
	t.Run("get_by_id", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &transaction.CreateRequest{
			Amount:    1,
			PartnerId: PartnerId,
			TxRef:     "test",
		}
		id, err := container.TransactionService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := container.TransactionService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		if id != resp.Id {
			t.Errorf("Expected id: %v Got: %v", id, resp.Id)
		}

	})

	t.Run("get_all", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &transaction.CreateRequest{
			Amount:    1,
			PartnerId: PartnerId,
			TxRef:     "test",
		}
		id, err := container.TransactionService.Create(ctx, in)
		print(id)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := container.TransactionService.GetAll(ctx)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected length: %v Want: %v", wantLen, len(resp.List))
		}
	})

	t.Run("get_by_date", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &transaction.CreateRequest{
			Amount:    1,
			PartnerId: PartnerId,
			TxRef:     "test",
		}
		id, err := container.TransactionService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := container.TransactionService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		resp_param, err := container.TransactionService.GetByParam(ctx, &transaction.GetByParamRequest{
			Date: resp.Date,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		wantLen := 1
		if wantLen != len(resp_param.List) {
			t.Errorf("Expected length: %v Want: %v", wantLen, len(resp_param.List))
		}
	})

	t.Run("get_by_partner_id", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &transaction.CreateRequest{
			Amount:    1,
			PartnerId: PartnerId,
			TxRef:     "test",
		}
		id, err := container.TransactionService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := container.TransactionService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		resp_param, err := container.TransactionService.GetByParam(ctx, &transaction.GetByParamRequest{
			PartnerId: resp.PartnerId,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		wantLen := 1
		if wantLen != len(resp_param.List) {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp_param.List))
		}
	})
}
