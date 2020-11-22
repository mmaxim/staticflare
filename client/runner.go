package client

import (
	"time"

	"mmaxim.org/staticflare/common"
)

type Runner struct {
	*common.DebugLabeler
	name, domain   string
	remoteIPSource RemoteIPSource
	dnsProvider    DNSProvider
}

func NewRunner(name, domain string, remoteIPSource RemoteIPSource, dnsProvider DNSProvider) *Runner {
	return &Runner{
		DebugLabeler:   common.NewDebugLabeler("Runner"),
		name:           name,
		domain:         domain,
		remoteIPSource: remoteIPSource,
		dnsProvider:    dnsProvider,
	}
}

func (r *Runner) Run() {
	lastIP, err := r.dnsProvider.GetDNS(r.name, r.domain)
	if err != nil {
		r.Debug("Run: failed to get initial IP: %s", err)
		return
	}
	for {
		ip, err := r.remoteIPSource.GetRemoteIP()
		if err != nil {
			r.Debug("Run: failed to get IP: %s", err)
		} else {
			r.Debug("Run: ip: %s lastIP: %s", ip, lastIP)
		}
		time.Sleep(time.Second)
	}
}
