package client

type DNSProvider interface {
	GetDNS(name, domain string) (string, error)
	SetDNS(name, domain, ip string) error
}

type RemoteIPSource interface {
	GetRemoteIP() (string, error)
}
