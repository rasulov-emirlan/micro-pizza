package psql_test

import (
	"database/sql"
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
	testHost     = "0.0.0.0"
	testPassword = "password"
	testDbName   = "micro-pizzas-users"
)

func setup() *dockertest.Resource {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not setup dockertest: %s", err.Error())
	}

	resource, err := pool.Run(
		"postgres", "14",
		[]string{
			fmt.Sprintf("POSTGRES_USER=%s", testUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", testPassword),
			fmt.Sprintf("POSTGRES_DB=%s", testDbName)})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}
	testPort = resource.GetPort("5432/tcp")

	if err := pool.Retry(func() error {
		db, err := sql.Open("postgres",
			fmt.Sprintf(
				"postgres://%s:%s@%s:%s/%s",
				testUser,
				testPassword,
				testHost,
				testPort,
				testDbName,
			))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("could not connect to database: %s", err.Error())
	}
	return resource
}

func TestMain(t *testing.M) {
	setup()
	code := t.Run()
	os.Exit(code)
}
