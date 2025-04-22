package main

import (
	"context"
	"flag"
	"log"
	"os"

	"mmaxim.org/staticflare/client"
	"mmaxim.org/staticflare/common"
)

type Options struct {
	Name              string
	Domain            string
	RemoteIPSourceURL string
	CloudFlareEmail   string
	CloudFlareAPIKey  string
	StathatEZKey      string
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
	flag.Usage()
	os.Exit(3)
}

func config() (opts Options) {
	flag.StringVar(&opts.Name, "name", os.Getenv("STATICFLARE_NAME"),
		"(required) subdomain to set on the domain (STATICFLARE_NAME env)")
	flag.StringVar(&opts.Domain, "domain", os.Getenv("STATICFLARE_DOMAIN"),
		"(required) DNS domain (STATICFLARE_DOMAIN env)")
	flag.StringVar(&opts.RemoteIPSourceURL, "ipurl", os.Getenv("STATICFLARE_IPURL"),
		"(required) URL for getting WAN IP from staticflared (STATICFLARE_DOMAIN env)")
	flag.StringVar(&opts.CloudFlareEmail, "cfemail", os.Getenv("STATICFLARE_CFEMAIL"),
		"(required) CloudFlare account email (STATICFLARE_CFEMAIL env)")
	flag.StringVar(&opts.CloudFlareAPIKey, "cfapikey", os.Getenv("STATICFLARE_CFAPIKEY"),
		"(required) CloudFlare API Key (STATICFLARE_CFAPIKEY env)")
	flag.StringVar(&opts.StathatEZKey, "stathatezkey", os.Getenv("STATICFLARE_STATHATEZKEY"),
		"(optional) StatHat EZ Key (STATICFLARE_STATHATEZKEY env)")
	flag.Parse()
	opts.check()
	return opts
}

func main() {
	opts := config()
	fullDNS := opts.Name + "." + opts.Domain

	// setup stats
	var statsProvider common.StatsProvider
	statsProvider = common.NewDummyStatsProvider()
	if opts.StathatEZKey != "" {
		log.Printf("StatHat EZ key provided\n")
		statsProvider = common.NewStathatStatsProvider("staticflared - "+fullDNS, opts.StathatEZKey)
	}

	// set up staticflared interface
	remoteIPSource := client.NewHTTPRemoteSource(opts.RemoteIPSourceURL,
		client.NewStaticFlaredHandler(), statsProvider)
	dnsProvider := client.NewCloudFlareDNSProvider()

	// set up CF
	if err := dnsProvider.Init(context.Background(), opts.CloudFlareEmail, opts.CloudFlareAPIKey); err != nil {
		log.Fatalf("failed to Init CloudFlare DNS provider: %s\n", err)
	}

	runner := client.NewRunner(opts.Name, opts.Domain, remoteIPSource, dnsProvider, statsProvider)
	runner.Run()
}
