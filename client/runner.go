package client

import (
	"context"
	"time"

	"mmaxim.org/staticflare/common"
)

type Runner struct {
	*common.DebugLabeler
	name, domain   string
	remoteIPSource RemoteIPSource
	dnsProvider    DNSProvider
	stats          common.StatsProvider
}

func NewRunner(name, domain string, remoteIPSource RemoteIPSource, dnsProvider DNSProvider,
	stats common.StatsProvider) *Runner {
	return &Runner{
		DebugLabeler:   common.NewDebugLabeler("Runner"),
		name:           name,
		domain:         domain,
		remoteIPSource: remoteIPSource,
		dnsProvider:    dnsProvider,
		stats:          stats.SetPrefix("Runner"),
	}
}

func (r *Runner) runOnce(ctx context.Context, lastIP string) (res string) {
	r.stats.CountOne("runOnce")
	ip, err := r.remoteIPSource.GetRemoteIP()
	if err != nil {
		r.stats.CountOne("runOnce - error")
		r.Debug("runOnce: failed to get IP: %s", err)
		return lastIP
	}
	if lastIP != ip {
		r.stats.CountOne("runOnce - update")
		r.Debug("runOnce: ip update: %s.%s: %s -> %s", r.name, r.domain, lastIP, ip)
		if err := r.dnsProvider.SetDNS(ctx, r.name, r.domain, ip); err != nil {
			r.stats.CountOne("runOnce - dns error")
			r.Debug("runOnce: failed to update: %s", err)
			return lastIP
		}
		return ip
	}
	return lastIP
}

func (r *Runner) Run() {
	ctx := context.Background()
	lastIP, err := r.dnsProvider.GetDNS(ctx, r.name, r.domain)
	if err != nil {
		r.Debug("Run: failed to get initial IP: %s", err)
		return
	}
	for {
		lastIP = r.runOnce(ctx, lastIP)
		time.Sleep(time.Second)
	}
}
