package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a users service"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
