//go:build local
// +build local

package main

import (
	"log"
	"net/http"

	handler "github.com/pranesh-2005/skill-icons-fast/api"
)

func main() {
	http.HandleFunc("/api/", handler.Handler)
	log.Println("Listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
