package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var port = "5012"

func hello(w http.ResponseWriter, r *http.Request) {
	var name, _ = os.Hostname()
	fmt.Fprintf(w, "Happy Coding!\n This request was processed by host: %s", name)
}

func main() {
	http.HandleFunc("/", hello)
	// get port env var
	portEnv := os.Getenv("APP_PORT")
	if len(portEnv) > 0 {
		port = portEnv
	}
	log.Printf("Listening on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}