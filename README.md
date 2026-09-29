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

## Building with Docker (xcaddy)

Caddy doesn't ship with third-party DNS modules built in; they're compiled in using [xcaddy](https://github.com/caddyserver/xcaddy). The [official Caddy Docker image](https://hub.docker.com/_/caddy) has a `:builder` variant with xcaddy preinstalled, which is the easiest way to build a custom Caddy binary with this module in a multi-stage `Dockerfile`:

```dockerfile
FROM caddy:builder AS builder

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    xcaddy build \
    --with github.com/caddy-dns/netcup

FROM caddy

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
```

The builder stage compiles a `caddy` binary with this module linked in; the final stage copies just that binary into a plain `caddy` image, keeping the final image small. You can add more `--with` lines for other Caddy modules (other DNS providers, plugins, etc.) in the same `xcaddy build` call.

### Using CloudDNS (REST API): pointing xcaddy at a patched libdns/netcup

This module depends on [libdns/netcup](https://github.com/libdns/netcup) for the actual netcup API calls (see [What does caddy-dns do, and what's libdns for?](#what-does-this-module-do-vs-libdnsnetcup) below). At the time of writing, the upstream `libdns/netcup` only supports netcup's legacy DNS API, not the newer CloudDNS REST API — see [libdns/netcup#4](https://github.com/libdns/netcup/issues/4). If your domain uses CloudDNS and upstream support isn't merged yet, you can build against a patched fork instead, without changing anything else in this module or your Caddyfile:

```dockerfile
FROM caddy:builder AS builder

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    xcaddy build \
    --with github.com/caddy-dns/netcup \
    --with github.com/libdns/netcup

FROM caddy

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
```

The second `--with module=replacement` line is an xcaddy [replace directive](https://github.com/caddyserver/xcaddy#in-line-replacements): it tells the Go build to substitute `github.com/YOUR-USER/libdns-netcup` (your fork, or any fork with CloudDNS support) wherever `github.com/caddy-dns/netcup` imports `github.com/libdns/netcup`. This module's own code and Caddyfile syntax stay exactly the same either way — replace `YOUR-USER` with your GitHub username (or point it at any other compatible fork/branch) once you have one published, then configure the `netcup` provider as usual:

```
bw.your.domain.com {
	tls {
		dns netcup {
			api_key {env.NETCUP_API_KEY}
		}
	}
}
```

leaving out `customer_number` and `api_password` is what selects the REST/CloudDNS code path in a patched `libdns/netcup`; see [Which API do I need?](#which-api-do-i-need) above.

## What does this module do vs. libdns/netcup?

This repository (`caddy-dns/netcup`) is a thin adapter: it registers a `netcup` DNS provider module with Caddy and translates Caddyfile/JSON config into a [`libdns/netcup`](https://github.com/libdns/netcup) `Provider` struct. It contains no netcup-specific API logic itself.

[`libdns/netcup`](https://github.com/libdns/netcup) is where the actual netcup API calls live (logging in, reading/writing DNS records), implemented against the generic [libdns](https://github.com/libdns/libdns) interfaces so the same code can be reused by other Go programs, not just Caddy. If you ever need to fix or extend netcup-specific behavior (like CloudDNS support above), that's the repository to change — this one (`caddy-dns/netcup`) only needs changes for Caddyfile/JSON syntax.
