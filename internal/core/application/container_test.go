package core

import (
	"database/sql"
	"os"
	"reflect"
	"testing"

	"b2b.nati011.github.com/config"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DependencyCheck struct {
	Service      string
	Dependencies []string
}

var serviceDependencies = []DependencyCheck{
	{
		Service:      "TemplateService",
		Dependencies: []string{},
	},
	{
		Service:      "EmailService",
		Dependencies: []string{"TemplateService"},
	},
	{
		Service:      "RoleService",
		Dependencies: []string{"ResourceService"},
	},
	{
		Service:      "AuthService",
		Dependencies: []string{"EmailService", "RoleService"},
	},
	{
		Service: "AuthMiddleware",
		Dependencies: []string{
			"AuthService",
			"RoleService",
			"ResourceService",
			"UserService",
		},
	},
	{
		Service:      "UserService",
		Dependencies: []string{"RoleService", "AuthService"},
	},
	{
		Service: "CheckoutService",
		Dependencies: []string{
			"PaymentService",
			"PaymentPartnerService",
			"TransactionService",
		},
	},
}

var db_pool *sql.DB
var application_container Container

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	db_pool = db_test_container.Setup()
	var cfg config.Config
	application_container = *NewContainer(
		db_pool,
		&cfg,
		config.DefaultPaginationBuilder().Build())
}

func teardown() {
	if db_pool != nil {
		db_pool.Close()
	}
}

func Test_ServiceInitializationOrder(t *testing.T) {
	t.Parallel()

	// Helper function to get service from container
	getService := func(serviceName string) interface{} {
		v := reflect.ValueOf(application_container)
		return reflect.Indirect(v).FieldByName(serviceName).Interface()
	}

	// Test each dependency relationship
	for _, dep := range serviceDependencies {
		t.Run(dep.Service, func(t *testing.T) {
			// Check if the main service exists and is initialized
			service := getService(dep.Service)
			require.NotNil(t, service, "Service %s should be initialized", dep.Service)

			// Check if all dependencies exist and are initialized
			for _, depName := range dep.Dependencies {
				dependency := getService(depName)
				require.NotNil(t, dependency,
					"%s requires %s to be initialized first",
					dep.Service, depName)
			}
		})
	}
}

func Test_ContainerValidity(t *testing.T) {
	t.Parallel()

	v := reflect.ValueOf(application_container)
	v = reflect.Indirect(v)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name

		if field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface {
			t.Run(fieldName, func(t *testing.T) {
				assert.False(t, field.IsNil(),
					"Service %s should not be nil", fieldName)
			})
		}
	}
}
