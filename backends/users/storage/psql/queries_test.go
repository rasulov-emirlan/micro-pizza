package psql_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/storage/psql"
)

var testPort string

var repo *psql.Repository

const (
	testUser     = "postgres"
	testPassword = "password"
	testHost     = "localhost"
	testDbName   = "micro-pizzas-users"
)

func setup() *dockertest.Resource {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not setup dockertest: %s", err.Error())
	}

	resource, err := pool.Run("postgres", "14", []string{fmt.Sprintf("POSTGRES_PASSWORD=%s", testPassword), fmt.Sprintf("POSTGRES_DB=%s", testDbName)})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}
	testPort = resource.GetPort("5432/tcp") // Set port used to communicate with Postgres
	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		r, err := psql.NewRepository(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", testUser, testPassword, testHost, testPort, testDbName), true)
		if err != nil {
			return err
		}
		r.Close()
		return err
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err.Error())
	}
	repo, err = psql.NewRepository(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", testUser, testPassword, testHost, testPort, testDbName), true)
	if err != nil {
		log.Fatalf("could not connect to db: %s", err.Error())
	}
	return resource
}

func TestMain(t *testing.M) {
	setup()
	code := t.Run()
	os.Exit(code)
}
