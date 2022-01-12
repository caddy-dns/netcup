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

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// netcup {
//     customer_number <customer_number>
//     api_key <api_key>
//     api_password <api_password>
// }
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
	if p.Provider.CustomerNumber == "" {
		return d.Err("missing customer number")
	}
	if p.Provider.APIKey == "" {
		return d.Err("missing API key")
	}
	if p.Provider.APIPassword == "" {
		return d.Err("missing API password")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
