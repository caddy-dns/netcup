# [netcup](https://www.netcup.de/) DNS module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records with the [netcup DNS API](https://ccp.netcup.net/run/webservice/servers/endpoint.php) using [libdns-netcup](https://github.com/libdns/netcup).

## Caddy module name

```
dns.providers.netcup
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) with your netcup credentials ([guide](https://www.netcup-wiki.de/wiki/CCP_API)) like so:

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
## Attention: Slow Netcup propagation time

NOTE: You may need to set an unexpectedly high propagation time (≥ 900 seconds) to give the netcup DNS time to propagate the entries! This may be annoying when calling certbot manually but should not be a problem in automated setups. In exceptional cases, 20 minutes may be required. See [coldfix/certbot-dns-netcup/issues/28](https://github.com/coldfix/certbot-dns-netcup/issues/28).

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
