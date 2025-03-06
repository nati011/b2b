package role

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	resource_db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/role/db"
	resource "b2b.nati011.github.com/internal/core/application/resource"
	role "b2b.nati011.github.com/internal/core/application/role"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testContainer role.TestContainer
var service role.Provider
var pgContainer *postgres.PostgresContainer
var db *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	var err error
	ctx := context.Background()

	pgContainer, err = RunContainer(ctx)
	if err != nil {
		panic(err)
	}

	connectionString, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	db, err = sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	service = role.NewRole(
		db_adapter.NewPostgres(
			db,
		), resource.NewResource(
			resource_db_adapter.NewPostgres(
				db,
			),
		),
	)

	testContainer = role.NewIntegrationTestContainer(db)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// ddl
	err = runMigration(db, "/home/natanel/personal/b2b_clean/b2b/migration/core_db.sql")
	if err != nil {
		log.Fatalf("Error running migration: %v", err)
	}

	// functions
	err = runMigration(db, "/home/natanel/personal/b2b_clean/b2b/migration/core_db_functions.sql")
	if err != nil {
		log.Fatalf("Error running migration: %v", err)
	}

}

func teardown() {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("could not begin transaction: %v", err)
	}

	// Get all table names
	var tables []string
	rows, err := tx.Query("SELECT tablename FROM pg_tables WHERE schemaname = 'public';")
	if err != nil {
		tx.Rollback()
		log.Fatalf("could not fetch table names: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			tx.Rollback()
			log.Fatalf("could not scan table name: %v", err)
		}
		tables = append(tables, table)
	}

	// Prepare the TRUNCATE statement
	if len(tables) > 0 {
		truncateQuery := "TRUNCATE TABLE " + strings.Join(tables, ", ") + " RESTART IDENTITY CASCADE;"
		_, err = tx.Exec(truncateQuery)
		if err != nil {
			tx.Rollback()
			log.Fatalf("could not truncate tables: %v", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("could not commit transaction: %v", err)
	}
}

func runMigration(db *sql.DB, filename string) error {
	// Read the SQL file
	sqlBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	// Execute the SQL
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("could not execute SQL: %w", err)
	}

	return nil
}

func RunContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	return postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),

		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
}

func Test_Timeout(t *testing.T) {

}

func Test_create_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	in := role.CreateRequest{
		Name: "test",
		Desc: "test",
	}

	_, err := service.Create(ctx, &in)
	if err != nil {
		t.Errorf("Failed to create role err: %v", err)
	}

	//get resource
	getResp, _ := service.GetAll(ctx)
	if len(getResp.List) == 0 {
		t.Errorf("No resources were created for role")
	}
}

func Test_create_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("duplicateName", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		in := role.CreateRequest{
			Desc: "test",
			Name: "test",
		}
		_, err := service.Create(ctx, &in)
		if err != nil {
			t.Errorf("Failed to create role err: %v", err)
		}

		//create duplicate
		wantErr := role.ErrDuplicateName
		got, err := service.Create(ctx, &in)
		if err != wantErr {
			switch err {
			case nil:
				t.Errorf("Expected err: %v Got err: %v", wantErr, err)
			default:
				t.Errorf("Failed to create role err: %v", err)
			}
		}
		if got != 0 {
			t.Errorf("Expected: %v, Got: %v", 0, got)
		}
	})

	t.Run("emptyName", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := role.CreateRequest{
			Desc: "test",
			Name: "",
		}

		wantErr := role.ErrEmptyName
		got, err := service.Create(ctx, &in)
		if err != wantErr {
			switch err {
			case nil:
				t.Errorf("Expected err: %v Got err: %v", wantErr, err)
			default:
				t.Errorf("Failed to create role err: %v", err)
			}
		}
		if got != 0 {
			t.Errorf("Expected: %v, Got: %v", 0, got)
		}
	})
}

func Test_delete_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	//setup
	ctx := context.Background()
	in := role.CreateRequest{
		Desc: "test",
		Name: "test",
	}
	id, _ := service.Create(ctx, &in)

	//delete resource
	err := service.Delete(ctx, id)
	if err != nil {
		t.Errorf("Failed to delete role err: %v", err)
	}

	//verify deletion
	resp, _ := service.Get(ctx, &role.GetRequest{
		Id:   id,
		Name: "",
	})
	if resp.Id == id {
		t.Error("Failed to delete role")
	}
}

func Test_delete_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		err := service.Delete(ctx, 1010)
		wantErr := role.ErrIdNotFound
		if err != wantErr {
			switch err {
			default:
				t.Errorf("Expected err: %q Got err: %q", wantErr, err)
			}
		}
	})
}

func Test_update_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	//setup
	ctx := context.Background()
	id, _ := service.Create(ctx, &role.CreateRequest{
		Desc: "test",
		Name: "test",
	})

	//update
	in := role.UpdateRequest{
		Id:   id,
		Desc: "test",
		Name: "tests",
	}
	resp, err := service.Update(ctx, &in)
	if err != nil {
		t.Errorf("Failed to update err %v", err)
	}
	if id != resp {
		t.Errorf("Failed to update role")
	}
}

func Test_update_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := role.UpdateRequest{
			Id:   1,
			Desc: "test",
			Name: "test",
		}
		wantErr := role.ErrIdNotFound
		resp, err := service.Update(ctx, &in)
		if err != wantErr {
			switch err {
			default:
				t.Errorf("Expected %v Got %v", wantErr, err)
			}
		}
		if resp != 0 {
			t.Errorf("Failed to update role")
		}
	})

	t.Run("duplicateName", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
		ctx := context.Background()

		//create resource with taken name
		service.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "taken",
		})

		// create resource
		id, _ := service.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})
		// update
		in := role.UpdateRequest{
			Id:   id,
			Desc: "test",
			Name: "taken",
		}
		wantErr := role.ErrDuplicateName
		got, err := service.Update(ctx, &in)
		if err != wantErr {
			switch err {
			default:
				t.Errorf("Expected %v Got %v", wantErr, err)
			}
		}
		if got != 0 {
			t.Errorf("Failed to update role")
		}

	})

	t.Run("emptyName", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})
		// update
		in := role.UpdateRequest{
			Id:   id,
			Desc: "tets",
			Name: "",
		}
		wantErr := role.ErrDuplicateName
		resp, err := service.Update(ctx, &in)
		if err != nil {
			switch err {
			case wantErr:
				t.Errorf("Expected %v Got %v", wantErr, err)
			default:
				t.Errorf("Failed to update err %v", err)
			}
		}
		if id != resp {
			t.Errorf("Failed to update role")
		}
	})

}

func Test_getRole_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("getById", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})

		// Get by Id
		in := role.GetRequest{
			Id: id,
		}
		got, err := service.Get(ctx, &in)
		if err != nil {
			t.Errorf("Failed to get role by Id err %v", err)
		}
		if got.Id != id {
			t.Errorf("Failed to get role by id")
		}
	})

	t.Run("getByName", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})

		// Get by Name
		in := role.GetRequest{
			Name: "test",
		}
		got, err := service.Get(ctx, &in)
		if err != nil {
			switch err {
			case role.ErrIdNotFound:
			default:
				t.Errorf("Failed to get role by Id err %v", err)
			}
		}
		if got.Id != id {
			t.Errorf("Failed to get role by id")
		}
	})
}

func Test_getRole_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("getById_notfound", func(t *testing.T) {
		//setup
		ctx := context.Background()

		// Get by Id
		in := role.GetRequest{
			Id: rand.Int(),
		}
		wantErr := role.ErrEmptyGetContent
		got, err := service.Get(ctx, &in)
		if err != wantErr {
			switch err {
			case wantErr:
			default:
				t.Errorf("Expected err: %v Got err %v", wantErr, err)
			}

		}
		if got.Id != 0 {
			t.Errorf("Failed, non existing resource found")
		}
	})

	t.Run("getByName", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		// Get by Name
		in := role.GetRequest{
			Name: "test",
		}
		wantErr := role.ErrEmptyGetContent
		got, err := service.Get(ctx, &in)
		if err != wantErr {
			t.Errorf("Expected err: %v Got err %v", wantErr, err)
		}
		if got.Id != 0 {
			t.Errorf("Failed, non existing resource found")
		}
	})
}

func Test_getAllResources_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	//setup
	ctx := context.Background()
	id, _ := service.Create(ctx, &role.CreateRequest{
		Desc: "test",
		Name: "test",
	})

	// Get All
	got, err := service.GetAll(ctx)
	if err != nil {
		t.Errorf("Failed to get role by Id err %v", err)
	}
	if len(got.List) == 0 || got.List[0].Id != id {
		t.Errorf("Failed to get role by id")
	}
}

func Test_hasResource_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("has", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		roleId, _ := service.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})

		// create resource with taken name
		resId, err := testContainer.ResourceService.Create(ctx, &resource.CreateRequest{
			Action: "test",
			Name:   "taken",
		})
		if err != nil {
			t.Errorf("Failed")
		}

		err = service.AddResource(ctx, &role.AddResourceRequest{
			ResourceId: resId,
			RoleId:     roleId,
		})
		if err != nil {
			t.Errorf("Failed to add resource")
		}

		in := role.HasResourceRequest{
			ResourceId: resId,
			RoleId:     roleId,
		}
		got, _ := service.HasResource(ctx, &in)
		if got != true {
			t.Errorf("Expected: %v Got: %v", true, got)
		}
	})

}

func Test_hasResource_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("has-not", func(t *testing.T) {
		//setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})

		in := role.HasResourceRequest{
			ResourceId: rand.Intn(200),
			RoleId:     id,
		}
		got, err := service.HasResource(ctx, &in)
		if err != nil {
			t.Errorf("Failed to check role resources %v", err)
		}
		if got != false {
			t.Errorf("Expected: %v Got: %v", false, got)
		}
	})
}

func Test_addResource_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	//setup
	ctx := context.Background()
	// create resource with taken name
	resId, err := testContainer.ResourceService.Create(ctx, &resource.CreateRequest{
		Action: "test",
		Name:   "taken",
	})
	if err != nil {
		t.Errorf("Failed")
	}

	//create role
	in := role.CreateRequest{
		Name: "test2",
		Desc: "test",
	}
	roleId, err := service.Create(ctx, &in)
	if err != nil {
		t.Errorf("Failed to add resource to role %v", err)
	}

	//add role to resource
	req := role.AddResourceRequest{
		ResourceId: resId,
		RoleId:     roleId,
	}
	err = service.AddResource(ctx, &req)
	if err != nil {
		t.Errorf("Failed to add resource to role %v", err)
	}
	inRes := role.HasResourceRequest{
		ResourceId: req.ResourceId,
		RoleId:     req.RoleId,
	}
	hasResource, err := service.HasResource(ctx, &inRes)
	if err != nil {
		t.Errorf("Failed to add resource to role %v", err)
	}
	if !hasResource {
		t.Errorf("Failed to add resource to role")
	}
}

func Test_addResource_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("resourceNotFound", func(t *testing.T) {
		ctx := context.Background()
		//create role
		in := role.CreateRequest{
			Name: "test2",
			Desc: "test",
		}
		roleId, err := service.Create(ctx, &in)
		if err != nil {
			t.Errorf("Failed to add resource to role %v", err)
		}

		//add role to resource
		req := role.AddResourceRequest{
			ResourceId: rand.Int(),
			RoleId:     roleId,
		}
		err = service.AddResource(ctx, &req)
		wantErr := role.ErrResourceNotFound
		if err != role.ErrResourceNotFound {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_removeResource_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	// setup
	ctx := context.Background()

	// create resource
	resId, err := testContainer.ResourceService.Create(ctx, &resource.CreateRequest{
		Action: "test",
		Name:   "test",
	})
	if err != nil {
		t.Errorf("Failed")
	}

	// create role
	roleId, err := service.Create(ctx, &role.CreateRequest{
		Name: "test2",
		Desc: "test",
	})
	if err != nil {
		t.Errorf("Failed to add resource to role %v", err)
	}

	// add role to resource
	req := role.AddResourceRequest{
		ResourceId: resId,
		RoleId:     roleId,
	}
	err = service.AddResource(ctx, &req)
	if err != nil {
		t.Errorf("Failed to add resource to role %v", err)
	}

	// check if role has resource
	hasResource, err := service.HasResource(ctx, &role.HasResourceRequest{
		ResourceId: req.ResourceId,
		RoleId:     req.RoleId,
	})
	if err != nil {
		t.Errorf("Failed %v", err)
	}
	if !hasResource {
		t.Errorf("Failed")
	}

	//remove resource
	err = service.RemoveResource(ctx, &role.RemoveResourceRequest{
		ResourceId: req.ResourceId,
		RoleId:     req.RoleId,
	})
	if err != nil {
		t.Errorf("Failed to remove resource from role %v", err)
	}

	//recheck
	hasResource, err = service.HasResource(ctx, &role.HasResourceRequest{
		ResourceId: req.ResourceId,
		RoleId:     req.RoleId,
	})
	if err != nil {
		t.Errorf("Failed %v", err)
	}
	if hasResource {
		t.Errorf("Failed to remove resource from role")
	}
}

func Test_removeResource_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("resource_not_found", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
		ctx := context.Background()

		// create resource
		resId, err := testContainer.ResourceService.Create(ctx, &resource.CreateRequest{
			Action: "test",
			Name:   "test",
		})
		if err != nil {
			t.Errorf("Failed")
		}

		// create role
		roleId, err := service.Create(ctx, &role.CreateRequest{
			Name: "test2",
			Desc: "test",
		})
		if err != nil {
			t.Errorf("Failed to add resource to role %v", err)
		}

		// add role to resource
		req := role.AddResourceRequest{
			ResourceId: resId,
			RoleId:     roleId,
		}
		err = service.AddResource(ctx, &req)
		if err != nil {
			t.Errorf("Failed to add resource to role %v", err)
		}

		// check if role has resource
		hasResource, err := service.HasResource(ctx, &role.HasResourceRequest{
			ResourceId: req.ResourceId,
			RoleId:     req.RoleId,
		})
		if err != nil {
			t.Errorf("Failed %v", err)
		}
		if !hasResource {
			t.Errorf("Failed")
		}

		//remove resource
		err = service.RemoveResource(ctx, &role.RemoveResourceRequest{
			ResourceId: req.ResourceId,
			RoleId:     req.RoleId,
		})
		if err != nil {
			t.Errorf("Failed to remove resource from role %v", err)
		}

		//recheck
		hasResource, err = service.HasResource(ctx, &role.HasResourceRequest{
			ResourceId: req.ResourceId,
			RoleId:     req.RoleId,
		})
		if err != nil {
			t.Errorf("Failed %v", err)
		}
		if hasResource {
			t.Errorf("Failed to remove resource from role")
		}

		//remove again
		err = service.RemoveResource(ctx, &role.RemoveResourceRequest{
			ResourceId: req.ResourceId,
			RoleId:     req.RoleId,
		})
		if err != nil {
			switch err {
			case role.ErrResourceNotFound:
			default:
				t.Errorf("Failed to remove resource from role %v", err)
			}
		}
	})
}
