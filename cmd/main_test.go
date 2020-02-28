package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	hello(w, r)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	got := w.Body.String()
	want := "Happy Coding"
	if !strings.Contains(w.Body.String(), "Happy Coding") {
		t.Fatalf("\nOops 🔥\x1b[91m Failed asserting that\x1b[39m\n"+
			"✘got: %v\n\x1b[92m"+
			"want Contains: %v\x1b[39m", got, want)
	}
}

func TestMain(m *testing.M) {
	// Use a helper function to ensure we run shutdown()
	// before calling os.Exit()
	os.Exit(mainHelper(m))
}

func mainHelper(m *testing.M) int {
	return m.Run()
}
