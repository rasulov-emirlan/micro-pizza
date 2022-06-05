package psql_test

import (
	"database/sql"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func setup() *dockertest.Resource {
	path := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("postgres", "postgres"),
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not setup dockertest: %s", err.Error())
	}

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "14",
			Env: []string{
				"POSTGRES_PASSWORD=postgres",
				"POSTGRES_USER=postgres",
				"POSTGRES_DB=micro-pizzas-users",
			},
		}, func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}
	path.Host = resource.Container.NetworkSettings.IPAddress + ":" + resource.GetPort("5432/tcp")
	log.Println(path.String() + "/micro-pizzas-users")
	pool.MaxWait = 120 * time.Second
	if err := pool.Retry(func() error {
		db, err := sql.Open("postgres", path.String()+"/micro-pizzas-users")
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("could not connect to database: %s", err.Error())
	}

	// kill container in 60 seconds
	resource.Expire(60)
	return resource
}

func TestMain(t *testing.M) {
	container := setup()
	code := t.Run()
	if err := container.Close(); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}
