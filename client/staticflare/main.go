package main

import (
	"log"
	"os"

	"mmaxim.org/staticflare/client"
)

func main() {
	remoteIPSource := client.NewHTTPRemoteSource("http://localhost:8080/info",
		client.NewStaticFlaredHandler())
	dnsProvider := client.NewCloudFlareDNSProvider()
	name := "staticdebug"
	domain := "mmaxim.org"

	// set up CF
	if err := dnsProvider.Init("mike.maxim@gmail.com", os.Getenv("STATICFLARE_CF_APIKEY")); err != nil {
		log.Fatalf("failed to Init CloudFlare DNS provider: %s\n", err)
	}

	runner := client.NewRunner(name, domain, remoteIPSource, dnsProvider)
	runner.Run()
}
