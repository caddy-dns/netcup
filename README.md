# [netcup](https://www.netcup.de/) DNS module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records with netcup's DNS APIs using [libdns-netcup](https://github.com/libdns/netcup).

## Caddy module name

```
dns.providers.netcup
```

## Which API do I need?

netcup exposes two different DNS APIs, and this module picks the right one automatically based on which config fields you set:

- **Legacy API**, for domains using "classic" DNS (shown as the **DNS** tab when you open the domain in the CCP): set `customer_number`, `api_key` and `api_password`. Supports full record management.
- **REST API**, for domains using **CloudDNS** (shown as the **CloudDNS** tab in the CCP): set only `api_key` — the key from the CCP's **API Keys** section, *not* **Legacy API Keys** — and omit `api_password`. `customer_number` isn't used in this mode. This API can currently only manage `_acme-challenge` records, which is exactly what's needed for the ACME DNS challenge that this module is for.

You can check which one your domain uses in the CCP under **Domains → 🔍 next to the domain**. See the [netcup API guide](https://www.netcup-wiki.de/wiki/CCP_API) and [libdns/netcup's README](https://github.com/libdns/netcup#configuration) for more details.

## Config examples

### Legacy API (classic DNS)

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) with your netcup credentials like so:

```json
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "netcup",
        "customer_number": "{env.NETCUP_CUSTOMER_NUMBER}",
        "api_key": "{env.NETCUP_API_KEY}",
        "api_password": "{env.NETCUP_API_PASSWORD}"
      }
    }
  }
}
```

or with the Caddyfile:

```
your.domain.com {
	...
	tls {
		dns netcup {
			customer_number {env.NETCUP_CUSTOMER_NUMBER}
			api_key {env.NETCUP_API_KEY}
			api_password {env.NETCUP_API_PASSWORD}
		}
	}
	...
}
```

### REST API (CloudDNS)

For a domain on CloudDNS, omit `customer_number` and `api_password` and use the `api_key` from the CCP's "API Keys" section:

```
your.domain.com {
	...
	tls {
		dns netcup {
			api_key {env.NETCUP_API_KEY}
		}
	}
	...
}
```

The rest of this README (propagation timing, resolvers) applies to both APIs, though CloudDNS deployment is usually noticeably faster (seconds rather than minutes) than the legacy API's propagation time discussed below.

## Attention: Slow Netcup propagation time (legacy API)

NOTE: You may need to set an unexpectedly high propagation time (≥ 900 seconds) to give the netcup DNS time to propagate the entries! This may be annoying when executing `caddy run/start` manually but should not be a problem in automated setups. In exceptional cases, 20 minutes may be required. See [coldfix/certbot-dns-netcup/issues/28](https://github.com/coldfix/certbot-dns-netcup/issues/28).

Use this config, (related [Caddy Forum Post](https://caddy.community/t/unable-to-issue-a-ssl-cert-via-acme-dns-for-netcup/21807)):
```
your.domain.com {
	...
	tls {
		dns netcup {
			...
		}
		propagation_timeout 900s
		propagation_delay 600s
		resolvers 1.1.1.1
	}
	...
}
```
