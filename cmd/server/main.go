package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/tell", http.HandlerFunc(tellHandler))

	// start an http server
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 100 * time.Millisecond,
		ConnState: func(conn net.Conn, state http.ConnState) {
			logger.Info("connection state changed", zap.String("state", state.String()))
		},
	}

	// TCP listener
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// start the server
	term := make(chan struct{})
	go func() {
		defer close(term)
		_ = server.Serve(lis)
	}()

	// capture sig signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// server should shutdown within 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-sig:
		_ = server.Shutdown(ctx)
	case <-term:
		log.Fatalf("server terminated without signal")
	}

	// wait for server to shutdown, it should stop within 5 seconds
	select {
	case <-ctx.Done():
		log.Fatalf("server didn't shutdown within timeout")
	case <-term:
		log.Printf("server shutdown gracefully")
	}
}

func tellHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
