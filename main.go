package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func getExternalIP(ctx context.Context) (string, error) {
	res, err := http.Get("wtfismyip.com/text")
	if err != nil {
		return "", fmt.Errorf("unable to get data: %w", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read body: %w", err)
	}

	return string(resBody), nil
}

func main() {
	log.Printf("First Pass")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ip, err := getExternalIP(ctx)
	log.Printf("Got IP: %v, %v", ip, err)
}
