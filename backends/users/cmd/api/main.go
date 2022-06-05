package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/internal/config"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/internal/storage/psql"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a users service"))
	})
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("postgres://%s:%s@0.0.0.0:5432/%s?sslmode=disable",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.DBname)
	repo, err := psql.NewRepository(url, true)
	if err != nil {
		log.Println(errors.Cause(err))
		log.Fatal(err)
	}
	defer repo.Close()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
