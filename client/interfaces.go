package client

import "context"

type DNSProvider interface {
	GetDNS(ctx context.Context, name, domain string) (string, error)
	SetDNS(ctx context.Context, name, domain, ip string) error
}

type RemoteIPSource interface {
	GetRemoteIP() (string, error)
}
