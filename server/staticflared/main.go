package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"mmaxim.org/staticflare/server"
)

type Options struct {
	Addr string
}

func (o Options) check() {
	if len(o.Addr) == 0 {
		usage()
	}
}

func usage() {
	fmt.Printf("usage: staticflared <--addr addr>\n")
	os.Exit(3)
}

func config() (opts Options) {
	flag.StringVar(&opts.Addr, "addr", os.Getenv("STATICFLARED_ADDR"), "listen addr for the server")
	flag.Parse()
	opts.check()
	return opts
}

func main() {
	opts := config()
	if err := server.NewServer(opts.Addr).Run(); err != nil {
		log.Printf("error running server: %s\n", err)
	}
}
