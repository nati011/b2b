package partner

import (
	"context"
	"os"
	"testing"
	"time"

	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/application/user"
)

var testContainer transaction.TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = transaction.NewPackageIntegrationTestContainer()
}

func Test_Validate_PartnerId_Upon_Transaction_Create_happyPath(t *testing.T) {
	ctx := context.Background()
	//create partner
	id, err := testContainer.PartnerService.Create(ctx, &partner.CreateRequest{
		Name:    "test",
		Icon:    "test",
		BaseUrl: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create err:%v", err)
	}

	//create user
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := user.CreateRequest{
		FirstName:  "Natanel",
		LastName:   "Jemaneh",
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

	//create transaction
	_, err = testContainer.TransactionService.Create(ctx, &transaction.CreateRequest{
		User_Id:    user_id,
		Amount:     100,
		Partner_Id: id,
	})
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
}

func Test_Validate_PartnerId_Upon_Transaction_Create_unhappyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	//create user
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := user.CreateRequest{
		FirstName:  "Natnael",
		LastName:   "Jemaneh",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test33",
		DOB:        parsedTime,
		ExternalId: "123",
	}
	user_id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	//create transaction
	_, err = testContainer.TransactionService.Create(ctx, &transaction.CreateRequest{
		User_Id:    user_id,
		Amount:     100,
		Partner_Id: 99,
	})
	wantErr := transaction.ErrPartnerDoesNotExist
	if err != wantErr {
		t.Errorf("Expected err: %v Got : %v", wantErr, err)
	}
}
