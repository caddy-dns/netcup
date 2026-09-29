package netcup

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/netcup"
)

// Provider lets Caddy read and manipulate DNS records hosted by this DNS provider.
type Provider struct{ *netcup.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.netcup",
		New: func() caddy.Module { return &Provider{new(netcup.Provider)} },
	}
}

// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	replacer := caddy.NewReplacer()
	p.Provider.CustomerNumber = replacer.ReplaceAll(p.Provider.CustomerNumber, "")
	p.Provider.APIKey = replacer.ReplaceAll(p.Provider.APIKey, "")
	p.Provider.APIPassword = replacer.ReplaceAll(p.Provider.APIPassword, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens.
//
// netcup exposes two different DNS APIs, and which one this module uses is
// chosen automatically based on which fields are set, so the same "netcup"
// directive works for both:
//
//   - Legacy API, for domains using "classic" DNS (the "DNS" tab in the
//     CCP): set customer_number, api_key and api_password.
//
//     netcup {
//     customer_number <customer_number>
//     api_key <api_key>
//     api_password <api_password>
//     }
//
//   - REST API, for domains using CloudDNS (the "CloudDNS" tab in the CCP):
//     set only api_key (the key from the CCP's "API Keys" section, not
//     "Legacy API Keys"), and omit api_password. This API can only manage
//     "_acme-challenge" records, which is exactly what's needed for the
//     ACME dns-01 challenge. customer_number is not used in this mode.
//
//     netcup {
//     api_key <api_key>
//     }
//
// See https://github.com/libdns/netcup for details on both modes.
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "customer_number":
				if d.NextArg() {
					p.Provider.CustomerNumber = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "api_key":
				if d.NextArg() {
					p.Provider.APIKey = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "api_password":
				if d.NextArg() {
					p.Provider.APIPassword = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.APIKey == "" {
		return d.Err("missing API key")
	}
	// api_password and customer_number belong together (legacy API); if
	// either is set, both must be, so that a typo'd/forgotten field can't
	// silently switch this into REST-API mode (which supports far less).
	if p.Provider.APIPassword != "" && p.Provider.CustomerNumber == "" {
		return d.Err("missing customer number (required together with api_password to use the legacy netcup API)")
	}
	if p.Provider.CustomerNumber != "" && p.Provider.APIPassword == "" {
		return d.Err("missing API password (required together with customer_number to use the legacy netcup API; omit both to use the new REST API for CloudDNS domains instead)")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
