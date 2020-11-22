package main

import "mmaxim.org/staticflare/client"

func main() {
	remoteIPSource := client.NewHTTPRemoteSource("http://localhost:8080/info",
		client.NewStaticFlaredHandler())
	runner := client.NewRunner(remoteIPSource)
	runner.Run()
}
