package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/qclaogui/golang-api-server/version"

	"github.com/qclaogui/golang-api-server/internal/log"
)

var port = "5012"

func hello(w http.ResponseWriter, r *http.Request) {
	log.Info("request from: " + r.RemoteAddr)
	var name, _ = os.Hostname()
	fmt.Fprintf(w, "Happy Coding! v1\nThis request was processed by host: %s", name)
}

func main() {
	log.Info(fmt.Sprintf("Starting the service...[commit: %s, build time: %s, release: %s]",
		version.Commit, version.BuildTime, version.Release))

	http.HandleFunc("/", hello)
	// get port env var
	portEnv := os.Getenv("APP_PORT")
	if len(portEnv) > 0 {
		port = portEnv
	}
	log.Info("Listening on port:" + port)

	log.Fatal(http.ListenAndServe(":"+port, nil).Error())
}
