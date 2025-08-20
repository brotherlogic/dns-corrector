package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func getExternalIP(ctx context.Context) (string, error) {
	res, err := http.Get("http://wtfismyip.com/text")
	if err != nil {
		return "", fmt.Errorf("unable to get data: %w", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read body: %w", err)
	}

	return string(resBody), nil
}

func resolveIP(ctx context.Context, address string) (string, error) {
	ips, err := net.LookupIP(address)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		log.Printf("Resolved %v -> %v", address, ip.String())
	}

	return ips[0].To16().String(), nil
}

func main() {
	for {
		log.Printf("First Pass")

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		ip, err := getExternalIP(ctx)
		log.Printf("Got IP: %v, %v", ip, err)

		ip, err = resolveIP(ctx, "gramophile-grpc.brotherlogic-backend.com")
		log.Printf("Resolved: %v and %v", ip, err)

		time.Sleep(time.Minute)
	}
}
