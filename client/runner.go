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

func (r *Runner) runOnce(lastIP string) (res string) {
	ip, err := r.remoteIPSource.GetRemoteIP()
	if err != nil {
		r.Debug("runOnce: failed to get IP: %s", err)
		return lastIP
	}
	if lastIP != ip {
		r.Debug("runOnce: ip update: %s.%s: %s -> %s", r.name, r.domain, lastIP, ip)
		if err := r.dnsProvider.SetDNS(r.name, r.domain, ip); err != nil {
			r.Debug("runOnce: failed to update: %s", err)
			return lastIP
		}
		return ip
	}
	return lastIP
}

func (r *Runner) Run() {
	lastIP, err := r.dnsProvider.GetDNS(r.name, r.domain)
	if err != nil {
		r.Debug("Run: failed to get initial IP: %s", err)
		return
	}
	for {
		lastIP = r.runOnce(lastIP)
		time.Sleep(time.Second)
	}
}
