package payment_partner

import (
	"context"
	"os"
	"testing"
)

var container TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = NewIntegrationTestContainer()
}

func Test_Create_Payment_Option_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:             "test",
			Icon:             "test",
			Init_payment_url: "https://google.com",
		}

		id, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check
		resp, err := container.PartnerService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get err: %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.Id)
		}
	})

	t.Run("inactive_by_default", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:             "test",
			Icon:             "test",
			Init_payment_url: "https://google.com",
		}

		id, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check
		resp, err := container.PartnerService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get err: %v", err)
		}
		if resp.Status != INACTIVE_STATUS {
			t.Errorf("Expected status: %v Got: %v", INACTIVE_STATUS, resp.Status)
		}
	})
}

func Test_Create_Payment_Option_unhappyPath(t *testing.T) {
	t.Run("name_mandatory", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := &CreateRequest{
			Icon:             "test",
			Init_payment_url: "https://google.com",
		}

		_, err := container.PartnerService.Create(ctx, in)
		wantErr := ErrNameIsNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("icon_mandatory", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:             "test",
			Init_payment_url: "https://google.com",
		}

		_, err := container.PartnerService.Create(ctx, in)
		wantErr := ErrIconIsNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("init_payment_url_mandatory", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := &CreateRequest{
			Name: "test",
			Icon: "test",
		}

		_, err := container.PartnerService.Create(ctx, in)
		wantErr := ErrUrlIsNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Activate_Payment_Option_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	ctx := context.Background()
	in := &CreateRequest{
		Name:             "test",
		Icon:             "test",
		Init_payment_url: "test",
	}

	id, err := container.PartnerService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	//activate
	err = container.PartnerService.Activate(ctx, id)
	if err != nil {
		t.Fatalf("Failed to activate payment option err: %v", err)
	}

	got, err := container.PartnerService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}

	if got.Status != ACTIVE_STATUS {
		t.Errorf("Expected status: %v Got: %v", got.Status, ACTIVE_STATUS)
	}
}

func Test_Activate_Payment_Option_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		err := container.PartnerService.Activate(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})

	t.Run("alreadyActive", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:             "test",
			Icon:             "test",
			Init_payment_url: "test",
		}

		id, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		err = container.PartnerService.Activate(ctx, id)
		if err != nil {
			t.Fatalf("Failed to activate payment option err: %v", err)
		}

		//re-activate
		err = container.PartnerService.Activate(ctx, id)
		wantErr := ErrPaymentOptionaAlreadyActive
		if err != wantErr {
			t.Errorf("Expected err : %v GotL %v", wantErr, err)
		}
	})
}

func Test_Deactivate_Payment_Option_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	ctx := context.Background()
	in := &CreateRequest{
		Name:             "test",
		Icon:             "test",
		Init_payment_url: "test",
	}

	id, err := container.PartnerService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	err = container.PartnerService.Activate(ctx, id)
	if err != nil {
		t.Fatalf("Failed to activate payment option err: %v", err)
	}

	err = container.PartnerService.Deactivate(ctx, id)
	if err != nil {
		t.Fatalf("Failed to deactivate payment option err: %v", err)
	}

	got, err := container.PartnerService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}
	if got.Status != INACTIVE_STATUS {
		t.Errorf("Expected status: %v Got: %v", got.Status, INACTIVE_STATUS)
	}
}

func Test_Deactivate_Payment_Option_unhappyPath(t *testing.T) {
	t.Run("alreadyInactive", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:             "test",
			Icon:             "test",
			Init_payment_url: "test",
		}

		id, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		err = container.PartnerService.Deactivate(ctx, id)
		wantErr := ErrPaymentOptionAlreadyInactive
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("id_not_found", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		err := container.PartnerService.Deactivate(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_All_Payment_Options_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	ctx := context.Background()
	//setup
	in := &CreateRequest{
		Name:             "test",
		Icon:             "test",
		Init_payment_url: "test",
	}

	_, err := container.PartnerService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	resp, err := container.PartnerService.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to get all payment options err: %v", err)
	}
	wantLen := 1
	if len(resp.List) != wantLen {
		t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp.List))
	}
}

func Test_Get_All_Payment_Options_unhappyPath(t *testing.T) {
	t.Run("empty_get_content", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		_, err := container.PartnerService.GetAll(ctx)
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_Active_Payment_Options_happyPath(t *testing.T) {
	t.Run("get_only_actives", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Name:             "test",
			Icon:             "test",
			Init_payment_url: "test",
		}

		id, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		_, err = container.PartnerService.GetActive(ctx)
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}

		err = container.PartnerService.Activate(ctx, id)
		if err != nil {
			t.Fatalf("Failed to activate payment option err: %v", err)
		}

		resp, err := container.PartnerService.GetActive(ctx)
		if err != nil {
			t.Fatalf("Failed to get all payment options err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp.List))
		}

		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got id: %v", resp.List[0].Id, id)
		}
	})
}

func Test_Get_Active_Payment_Options_unhappyPath(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		_, err := container.PartnerService.GetActive(ctx)
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_Payment_Options_ByParam_happyPath(t *testing.T) {
	t.Run("get_by_name", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Name:             "test",
			Icon:             "test",
			Init_payment_url: "test",
		}

		id, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := container.PartnerService.GetByParam(ctx, &GetByParamRequest{
			Name: "test",
		})
		if err != nil {
			t.Fatalf("Failed to get by param err: %v", err)
		}
		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got id: %v", id, resp.List[0].Id)
		}
	})
}

func Test_Get_Payment_Options_ByParam_unhappyPath(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		_, err := container.PartnerService.GetByParam(ctx, &GetByParamRequest{})
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}
