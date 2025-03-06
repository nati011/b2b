package category

import (
	"context"
	"os"
	"testing"

	db "b2b.nati011.github.com/internal/adapter/secondary/domain/catalogue/category"
)

var service Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	service = NewCategory(
		db.NewMock(),
	)
}

func Test_add_category_happyPath(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name: "test",
			Desc: "test",
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create category err: %v", err)
		}

		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get category err: %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v Got id: %v", id, resp.Id)
		}
	})

}

func Test_add_category_unhappyPath(t *testing.T) {
	t.Run("name_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Desc: "test",
		}
		_, err := service.Create(ctx, in)
		wantErr := ErrNameIsNotSupplied
		if err != wantErr {
			t.Fatalf("Expected err: %v Want err: %v", wantErr, err)
		}
	})

	t.Run("desc_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name: "test",
		}
		_, err := service.Create(ctx, in)
		wantErr := ErrDescIsNotSupplied
		if err != wantErr {
			t.Fatalf("Expected err: %v Want err: %v", wantErr, err)
		}
	})

}

func Test_remove_category_happyPath(t *testing.T) {
	//setup
	ctx := context.Background()
	in := &CreateRequest{
		Name: "test",
		Desc: "test",
	}
	id, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create category err: %v", err)
	}

	//remove
	err = service.Remove(ctx, id)
	if err != nil {
		t.Fatalf("Failed to remove category err: %v", err)
	}

	//check
	_, err = service.Get(ctx, id)
	wantErr := ErrIdNotFound
	if err != wantErr {
		t.Errorf("Expected err: %v Got err: %v", wantErr, err)
	}
}

func Test_remove_category_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		ctx := context.Background()
		err := service.Remove(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_get_happyPath(t *testing.T) {
	// setup
	ctx := context.Background()
	in := &CreateRequest{
		Name: "test",
		Desc: "test",
	}
	id, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create category err: %v", err)
	}

	// get
	resp, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get category err: %v", err)
	}

	if resp.Id != id {
		t.Errorf("Expected id: %v Got id: %v", resp.Id, id)
	}
}

func Test_get_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.Get(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_get_all_happyPath(t *testing.T) {
	// setup
	ctx := context.Background()
	in := &CreateRequest{
		Name: "test",
		Desc: "test",
	}
	_, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create category err: %v", err)
	}

	// get
	resp, err := service.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to get category err: %v", err)
	}

	wantLen := 1
	if len(resp.List) != wantLen {
		t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp.List))
	}
}

func Test_get_all_unhappyPath(t *testing.T) {
	// setup
	ctx := context.Background()
	_, err := service.GetAll(ctx)
	wantErr := ErrEmptyGetContent
	if err != wantErr {
		t.Errorf("Expected err: %v Got err: %v", wantErr, err)
	}
}
