package source

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURL_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))
	defer srv.Close()

	u := URL{Raw: srv.URL}
	rc, err := u.Open(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rc.Close() }()
	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hello" {
		t.Fatalf("got %q", b)
	}
}

func TestURL_StatusError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	u := URL{Raw: srv.URL}
	_, err := u.Open(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
