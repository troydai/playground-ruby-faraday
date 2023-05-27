package main

import (
	"io"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	client := &http.Client{}
	for i := 0; i < 3; i++ {
		resp, err := client.Get("http://localhost:8080/tell")
		if err != nil {
			logger.Fatal("failed GET request", zap.Error(err))
		}

		io.ReadAll(resp.Body)
		resp.Body.Close()
	}
}
