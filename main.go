package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Marineholmen-Makerspace/members-management/m4m"
)

func StatusHandler(w http.ResponseWriter, _ *http.Request) {
	cfg := m4m.GetConfig()
	_, _ = w.Write([]byte(fmt.Sprintf("Using FabMab account: %d", cfg.FabMan.Account)))
}

func main() {
	http.HandleFunc("/_wh/stripe", m4m.StripeWebHookHandler)
	http.HandleFunc("/status", StatusHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
