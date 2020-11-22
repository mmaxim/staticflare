package client

import (
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	"mmaxim.org/staticflare/common"
)

type CloudFlareDNSProvider struct {
	*common.DebugLabeler
	api     *cloudflare.API
	zoneIDs map[string]string
}

func NewCloudFlareDNSProvider() *CloudFlareDNSProvider {
	return &CloudFlareDNSProvider{
		DebugLabeler: common.NewDebugLabeler("CloudFlareDNSProvider"),
		zoneIDs:      make(map[string]string),
	}
}

func (c *CloudFlareDNSProvider) Init(email, token string) error {
	api, err := cloudflare.New(token, email)
	if err != nil {
		return err
	}
	c.api = api
	zones, err := c.api.ListZones()
	if err != nil {
		c.Debug("Init: failed to list zones: %s", err)
		return err
	}
	for _, zone := range zones {
		c.zoneIDs[zone.Name] = zone.ID
		c.Debug("Init: id: %s -> %s", zone.Name, zone.ID)
	}
	return nil
}

func (c *CloudFlareDNSProvider) getZoneID(domain string) (string, error) {
	zoneID, ok := c.zoneIDs[domain]
	if !ok {
		return zoneID, fmt.Errorf("unknown domain: %s", domain)
	}
	return zoneID, nil
}

func (c *CloudFlareDNSProvider) getDNSRecord(name, domain string) (res cloudflare.DNSRecord, err error) {
	zoneID, err := c.getZoneID(domain)
	if err != nil {
		return res, err
	}
	c.Debug("getDNSRecord: using zoneID: %s", zoneID)
	recs, err := c.api.DNSRecords(zoneID, cloudflare.DNSRecord{})
	if err != nil {
		c.Debug("getDNSRecord: failed to get records: %s", err)
		return res, err
	}
	recname := name + "." + domain
	for _, rec := range recs {
		if rec.Name == recname {
			return rec, nil
		}
	}
	return res, fmt.Errorf("failed to get DNS record for name: %s domain: %s", name, domain)
}

func (c *CloudFlareDNSProvider) SetDNS(name, domain, ip string) error {
	rec, err := c.getDNSRecord(name, domain)
	if err != nil {
		c.Debug("SetDNS: failed to get DNS record: %s", err)
		return err
	}
	rec.Content = ip
	if err := c.api.UpdateDNSRecord(rec.ZoneID, rec.ID, rec); err != nil {
		c.Debug("SetDNS: failed to update: %s", err)
		return err
	}
	return nil
}

func (c *CloudFlareDNSProvider) GetDNS(name, domain string) (res string, err error) {
	rec, err := c.getDNSRecord(name, domain)
	if err != nil {
		return res, err
	}
	return rec.Content, nil
}
