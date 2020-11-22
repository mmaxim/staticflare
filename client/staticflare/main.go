package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"mmaxim.org/staticflare/client"
)

type Options struct {
	Name              string
	Domain            string
	RemoteIPSourceURL string
	CloudFlareEmail   string
	CloudFlareAPIKey  string
}

func (o Options) check() {
	if len(o.Name) == 0 ||
		len(o.Domain) == 0 ||
		len(o.RemoteIPSourceURL) == 0 ||
		len(o.CloudFlareEmail) == 0 ||
		len(o.CloudFlareAPIKey) == 0 {
		usage()
	}
}

func usage() {
	fmt.Printf("usage: staticflare <--name name_arg> <--domain domain_arg> <--ipurl ipurl_arg> <--cfemail cfemail_arg> <--cfapikey cfapikey_arg>\n")
	os.Exit(3)
}

func config() (opts Options) {
	flag.StringVar(&opts.Name, "name", os.Getenv("STATICFLARE_NAME"), "subdomain to set on the domain")
	flag.StringVar(&opts.Domain, "domain", os.Getenv("STATICFLARE_DOMAIN"), "DNS domain")
	flag.StringVar(&opts.RemoteIPSourceURL, "ipurl", os.Getenv("STATICFLARE_IPURL"),
		"URL for getting WAN IP")
	flag.StringVar(&opts.CloudFlareEmail, "cfemail", os.Getenv("STATICFLARE_CFEMAIL"),
		"CloudFlare account email")
	flag.StringVar(&opts.CloudFlareAPIKey, "cfapikey", os.Getenv("STATICFLARE_CFAPIKEY"),
		"CloudFlare API Key")
	flag.Parse()
	opts.check()
	return opts
}

func main() {
	opts := config()
	remoteIPSource := client.NewHTTPRemoteSource(opts.RemoteIPSourceURL,
		client.NewStaticFlaredHandler())
	dnsProvider := client.NewCloudFlareDNSProvider()

	// set up CF
	if err := dnsProvider.Init(opts.CloudFlareEmail, opts.CloudFlareAPIKey); err != nil {
		log.Fatalf("failed to Init CloudFlare DNS provider: %s\n", err)
	}

	runner := client.NewRunner(opts.Name, opts.Domain, remoteIPSource, dnsProvider)
	runner.Run()
}
