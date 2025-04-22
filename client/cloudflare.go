package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"mmaxim.org/staticflare/common"
)

type CloudFlareDNSProvider struct {
	*common.DebugLabeler
	api     *cloudflare.Client
	zoneIDs map[string]string
}

func NewCloudFlareDNSProvider() *CloudFlareDNSProvider {
	return &CloudFlareDNSProvider{
		DebugLabeler: common.NewDebugLabeler("CloudFlareDNSProvider"),
		zoneIDs:      make(map[string]string),
	}
}

func (c *CloudFlareDNSProvider) Init(ctx context.Context, email, token string) error {
	c.api = cloudflare.NewClient(option.WithAPIKey(token), option.WithAPIEmail(email))
	return nil
}

func (c *CloudFlareDNSProvider) getZoneID(ctx context.Context, domain string) (string, error) {
	c.Debug("CF: getZoneID: domain: %s", domain)
	zones, err := c.api.Zones.List(ctx, zones.ZoneListParams{
		Name: cloudflare.F(domain),
	})
	if err != nil {
		return "", err
	}
	if len(zones.Result) == 0 {
		return "", errors.New("not zone for domain")
	}
	return zones.Result[0].ID, nil
}

func (c *CloudFlareDNSProvider) SetDNS(ctx context.Context, name, domain, ip string) error {
	c.Debug("CF: SetDNS: setting DNS record: name: %s domain: %s ip: %s", name, domain, ip)
	rec, zoneID, err := c.getDNSRecord(ctx, name, domain)
	if err != nil {
		return err
	}
	updateParams := dns.RecordUpdateParams{
		ZoneID: cloudflare.F(zoneID),
		Record: dns.ARecordParam{
			Content: cloudflare.F(ip),
			Name:    cloudflare.F(name),
			Type:    cloudflare.F(dns.ARecordTypeA),
		},
	}
	c.Debug("CF: SetDNS: performing update")
	if _, err = c.api.DNS.Records.Update(ctx, rec.ID, updateParams); err != nil {
		return err
	}
	c.Debug("CF: SetDNS: success")
	return nil
}

func (c *CloudFlareDNSProvider) getDNSRecord(ctx context.Context, name, domain string) (res dns.RecordResponse, zoneID string, err error) {
	c.Debug("CF: getDNSRecord: fetching DNS record: name: %s domain: %s", name, domain)
	if zoneID, err = c.getZoneID(ctx, domain); err != nil {
		return res, zoneID, err
	}
	c.Debug("CF: getDNSRecord: zoneID: %s", zoneID)
	lres, err := c.api.DNS.Records.List(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(zoneID),
		Name: cloudflare.F(dns.RecordListParamsName{
			Exact: cloudflare.F(name + "." + domain),
		}),
	})
	if err != nil {
		return res, zoneID, err
	}
	if len(lres.Result) == 0 {
		return res, zoneID, fmt.Errorf("no records found: name: %s domain: %s", name, domain)
	}
	res = lres.Result[0]
	c.Debug("CF: getDNSRecord: found: content: %s", res.Content)
	return res, zoneID, nil
}

func (c *CloudFlareDNSProvider) GetDNS(ctx context.Context, name, domain string) (res string, err error) {
	c.Debug("CF: GetDNS: fetching DNS record: name: %s domain: %s", name, domain)
	rec, _, err := c.getDNSRecord(ctx, name, domain)
	if err != nil {
		return res, err
	}
	return rec.Content, nil
}
