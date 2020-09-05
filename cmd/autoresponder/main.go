package main

import (
	"fmt"
	"net/http"

	"github.com/jgengo/AutoResponder/internal/handler"
)

func main() {
	http.HandleFunc("/events", handler.EventHandler)
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":8080", nil)
}
