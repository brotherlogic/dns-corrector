package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/cloudflare/cloudflare-go/v5"
	"github.com/cloudflare/cloudflare-go/v5/dns"
	"github.com/cloudflare/cloudflare-go/v5/option"
)

type server struct {
	client *cloudflare.Client
}

func (s *server) runCorrection(ctx context.Context, newIP string) error {
	s.client.DNS.Records.Update(ctx, "monitoring", dns.RecordUpdateParams{
		ZoneID: cloudflare.F("brotherlogic-backend.com"),
	})
	return nil
}

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

func (s *server) loop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ip, err := getExternalIP(ctx)
	log.Printf("Got IP: %v, %v", ip, err)

	ip, err = resolveIP(ctx, "gramophile-grpc.brotherlogic-backend.com")
	log.Printf("Resolved: %v and %v", ip, err)
}

func (s *server) run() {

	for {
		s.loop()
		time.Sleep(time.Minute)
	}
}

func main() {

	client := cloudflare.NewClient(
		option.WithAPIToken("TOKEN"),
	)

	s := &server{
		client: client,
	}

	s.run()

}
