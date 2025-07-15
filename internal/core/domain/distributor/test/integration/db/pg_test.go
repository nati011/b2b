package db

import (
	"context"
	"database/sql"
	"log"
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

func Test_Timeout(t *testing.T) {
}

func Test_Read(t *testing.T) {
	t.Run("get", func(t *testing.T) {
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

		//check
		resp, err := testContainer.DistributorService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v Got :%v", id, resp.Id)
		}
	})

	t.Run("get_all", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		// setup
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

		resp, err := testContainer.DistributorService.GetAll(ctx)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(resp.List))
		}
		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.List[0].Id)
		}
	})

	t.Run("get_by_name", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		// setup
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
		resp, err := testContainer.DistributorService.GetByParam(ctx, &distributor.GetByParamRequest{
			Name: in.FirstName + in.LastName,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(resp.List))
		}
		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.List[0].Id)
		}
	})

	t.Run("get_by_tin", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		// setup
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
		resp, err := testContainer.DistributorService.GetByParam(ctx, &distributor.GetByParamRequest{
			Tin: in.Tin,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", len(resp.List), wantLen)
		}
		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.List[0].Id)
		}
	})

	t.Run("get_by_status", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		// setup
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
		err = testContainer.DistributorService.Activate(ctx, id)
		if err != nil {
			t.Fatalf("Failed to activate err: %v", err)
		}
		resp, err := testContainer.DistributorService.GetByParam(ctx, &distributor.GetByParamRequest{
			Status: "ACTIVE",
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", len(resp.List), wantLen)
		}
		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.List[0].Id)
		}
	})
	t.Run("getByUserId", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		// setup
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
		log.Printf("Id, %v", id)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		// distributor, err := testContainer.DistributorService.Get(ctx, id)
		// if err != nil {
		// 	t.Fatalf("Failed to create user %v", err)
		// }
		dist, err := testContainer.DistributorService.GetByUserId(ctx, id)
		if err != nil {
			t.Errorf("Failed to get distirbutor %v", err)
		}
		log.Printf("Distribtors %v", dist)
		if dist.Id != id {
			t.Fatalf("Want id %v got id %v", id, dist.Id)
		}

	})
	t.Run("get_all_user_agents", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

		_, err = testContainer.DistributorService.GetAllUsers(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get user agents %v", err)
		}

	})
}

func Test_Write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
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
			LicenceURL:  "example.com",

			Email: "test@gmail.com",
		}
		id, err := testContainer.DistributorService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check
		resp, err := testContainer.DistributorService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v Got :%v", id, resp.Id)
		}
	})

	t.Run("create_user_agent", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

			Email: "test@gmail.com",
		}
		id, err := testContainer.DistributorService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check if user agent has been created
		_, err = testContainer.DistributorService.GetAllUsers(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get user agents %v", err)
		}
	})

	t.Run("updateName", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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
		update_in := distributor.UpdateRequest{
			Id:   id,
			Name: "test",
		}
		_, err = testContainer.DistributorService.Update(ctx, &update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		//check
		resp, err := testContainer.DistributorService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if resp.Name != update_in.Name {
			t.Errorf("Expected name: %v Got:%v", update_in.Name, resp.Name)
		}
	})

	t.Run("updateTin", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
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
		update_in := distributor.UpdateRequest{
			Id:  id,
			Tin: "1234567891",
		}
		_, err = testContainer.DistributorService.Update(ctx, &update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		// check
		resp, err := testContainer.DistributorService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		if resp.Tin != update_in.Tin {
			t.Errorf("Expected Tin: %v Got:%v", update_in.Tin, resp.Tin)
		}
	})

	t.Run("activate", func(t *testing.T) {
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

		err = testContainer.DistributorService.Activate(ctx, id)
		if err != nil {
			t.Fatalf("Failed to activate err: %v", err)
		}

		// check
		resp, err := testContainer.DistributorService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		if resp.IsActive != true {
			t.Errorf("Expected isActive: %v Got:%v", true, resp.IsActive)
		}
	})

	t.Run("deactivate", func(t *testing.T) {
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

		err = testContainer.DistributorService.Activate(ctx, id)
		if err != nil {
			t.Fatalf("Failed to activate err: %v", err)
		}

		err = testContainer.DistributorService.Dectivate(ctx, id)
		if err != nil {
			t.Fatalf("Failed to activate err: %v", err)
		}

		// check
		resp, err := testContainer.DistributorService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		if resp.IsActive != false {
			t.Errorf("Expected isActive: %v Got:%v", false, resp.IsActive)
		}
	})
}
