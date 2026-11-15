package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"
)

func listen(t *testing.T) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	return ln
}

func TestRunServesUntilCancelled(t *testing.T) {
	ln := listen(t)
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
		ReadHeaderTimeout: time.Second,
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Run(ctx, srv, ln, time.Second) }()

	resp, err := http.Get("http://" + ln.Addr().String())
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d", resp.StatusCode)
	}

	cancel()
	if err := <-done; err != nil {
		t.Fatalf("Run returned %v, want nil", err)
	}
}

func TestRunReturnsServeError(t *testing.T) {
	ln := listen(t)
	_ = ln.Close()
	srv := &http.Server{ReadHeaderTimeout: time.Second}

	err := Run(context.Background(), srv, ln, time.Second)
	if err == nil {
		t.Fatal("expected an error from a closed listener")
	}
}

func TestRunReportsShutdownTimeout(t *testing.T) {
	ln := listen(t)
	started := make(chan struct{})
	release := make(chan struct{})
	srv := &http.Server{
		Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			close(started)
			<-release
		}),
		ReadHeaderTimeout: time.Second,
	}
	defer close(release)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Run(ctx, srv, ln, 10*time.Millisecond) }()
	go func() {
		resp, err := http.Get("http://" + ln.Addr().String())
		if err == nil {
			_ = resp.Body.Close()
		}
	}()

	<-started
	cancel()
	if err := <-done; !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Run returned %v, want deadline exceeded", err)
	}
}
