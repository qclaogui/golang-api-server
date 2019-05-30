package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/qclaogui/golang-api-server/version"

	"log"
)

var port = "5012"

var sourceLink = "https://github.com/qclaogui/golang-api-server"

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("request from: " + r.RemoteAddr)
	var ver = fmt.Sprintf("Build on %s [commit: %s, build time: %s, release: %s]", runtime.Version(),
		version.Commit, version.BuildTime, version.Release)

	var name, _ = os.Hostname()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<br/><center><h1>Happy Coding </h1><br/><code>%s</code><p><a href=%q target=_blank>source code</a></p></center><hr><br/>"+
		"<center>this request was processed by host: %s</center>", ver, sourceLink, name)
}

func main() {
	log.Println(fmt.Sprintf("Starting the service...[commit: %s, build time: %s, release: %s]",
		version.Commit, version.BuildTime, version.Release))

	http.HandleFunc("/", hello)
	// get port env var
	portEnv := os.Getenv("APP_PORT")
	if len(portEnv) > 0 {
		port = portEnv
	}
	log.Println("Listening on port:" + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
