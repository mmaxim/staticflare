package client

type DNSProvider interface {
	SetDNS(name, domain, ip string) error
}

type RemoteIPSource interface {
	GetRemoteIP() (string, error)
}
