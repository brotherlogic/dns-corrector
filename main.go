package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go/v5"
	"github.com/cloudflare/cloudflare-go/v5/dns"
	"github.com/cloudflare/cloudflare-go/v5/option"
)

type server struct {
	client *cloudflare.Client
}

func (s *server) runCorrection(ctx context.Context, newIP string) error {
	client := cloudflare.NewClient(
		option.WithAPIToken(os.Getenv("CLOUDFLARE_TOKEN")),
	)

	iter, err := client.DNS.Records.List(context.Background(), dns.RecordListParams{
		ZoneID: cloudflare.F("ee08022dbf5b9233d104a2b7a1778a82"),
	})
	if err != nil {
		log.Fatalf("Bad: %v", err)
	}

	for _, value := range iter.Result {
		_, err := client.DNS.Records.Edit(context.Background(), value.ID,
			dns.RecordEditParams{
				ZoneID: cloudflare.F("ee08022dbf5b9233d104a2b7a1778a82"),
				Body: &dns.RecordEditParamsBody{
					Name:    cloudflare.F(value.Name),
					Content: cloudflare.F(newIP),
				},
			})
		if err != nil {
			return err
		}
	}

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

	ipe, err := getExternalIP(ctx)
	log.Printf("Got IP: %v, %v", ipe, err)

	ipi, err := resolveIP(ctx, "gramophile-grpc.brotherlogic-backend.com")
	log.Printf("Resolved: %v and %v", ipi, err)

	if ipi != ipe {
		log.Printf("Running Correction")
		err = s.runCorrection(ctx, ipe)
		if err != nil {
			log.Fatalf("Unable to run correction: %v", err)
		}
	} else {
		log.Printf("No need for correction")
	}
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
