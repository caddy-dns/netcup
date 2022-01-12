# netcup DNS module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records with the [netcup DNS API](https://ccp.netcup.net/run/webservice/servers/endpoint.php).

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
