package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rasulov-emirlan/micro-pizzas/backends/users/config"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/migrations"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a users service"))
	})
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(20 * time.Second)
	if err := migrations.Up(
		fmt.Sprintf(
			"postgres://%s:%s@database:5432/%s?sslmode=disable",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.DBname),
	); err != nil {
		log.Fatal(err)
	}
	log.Println("just got up our migrations brother")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
