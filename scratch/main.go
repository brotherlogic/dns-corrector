package main

import (
	"context"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go/v5"
	"github.com/cloudflare/cloudflare-go/v5/dns"
	"github.com/cloudflare/cloudflare-go/v5/option"
)

func main() {
	newIP := "76.191.164.154"
	log.Printf("Running correction: %v", newIP)

	client := cloudflare.NewClient(
		option.WithAPIToken(os.Args[1]),
	)

	iter, err := client.DNS.Records.List(context.Background(), dns.RecordListParams{
		ZoneID: cloudflare.F("ee08022dbf5b9233d104a2b7a1778a82"),
	})
	if err != nil {
		log.Fatalf("Bad: %v", err)
	}

	for _, value := range iter.Result {
		log.Printf("HERE: %v -> %v", value.Name, value.ID)
	}

	val, err := client.DNS.Records.Edit(context.Background(), "0dbaac4ed89b37fce2443236fc8f3beb",
		dns.RecordEditParams{
			ZoneID: cloudflare.F("ee08022dbf5b9233d104a2b7a1778a82"),
			Body: &dns.RecordEditParamsBody{
				Name: cloudflare.F("registry.brotherlogic-backend.com"),
			},
		})

	log.Printf("Update: %v, %v", val, err)

}
