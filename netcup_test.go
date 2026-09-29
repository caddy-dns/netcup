package netcup

import (
	"fmt"
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/netcup"
)

func TestUnmarshalCaddyFileExtractsApiToken(t *testing.T) {
	tests := []string{
		`netcup {
			customer_number 12345
			api_key apikey
			api_password pw123
		}`}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			// given
			dispenser := caddyfile.NewTestDispenser(tc)
			p := Provider{&netcup.Provider{}}
			// when
			err := p.UnmarshalCaddyfile(dispenser)
			// then
			if err != nil {
				t.Errorf("UnmarshalCaddyfile failed with %v", err)
				return
			}

			expectedCustomerNumber := "12345"
			actualCustomerNumber := p.Provider.CustomerNumber
			if expectedCustomerNumber != actualCustomerNumber {
				t.Errorf("Expected CustomerNumber to be '%s' but got '%s'", expectedCustomerNumber, actualCustomerNumber)
			}

			expectedAPIKey := "apikey"
			actualAPIKey := p.Provider.APIKey
			if expectedAPIKey != actualAPIKey {
				t.Errorf("Expected CustomerNumber to be '%s' but got '%s'", expectedAPIKey, actualAPIKey)
			}

			expectedAPIPassword := "pw123"
			actualAPIPassword := p.Provider.APIPassword
			if expectedAPIPassword != actualAPIPassword {
				t.Errorf("Expected CustomerNumber to be '%s' but got '%s'", expectedAPIPassword, actualAPIPassword)
			}
		})
	}
}

// TestUnmarshalCaddyFileRestApiKeyOnly verifies that a Caddyfile block with
// only api_key (no api_password, no customer_number) parses successfully,
// which is what selects the new REST API for CloudDNS domains.
func TestUnmarshalCaddyFileRestApiKeyOnly(t *testing.T) {
	tc := `netcup {
		api_key resttestkey
	}`

	dispenser := caddyfile.NewTestDispenser(tc)
	p := Provider{&netcup.Provider{}}
	if err := p.UnmarshalCaddyfile(dispenser); err != nil {
		t.Fatalf("UnmarshalCaddyfile failed with %v", err)
	}

	if got := p.Provider.APIKey; got != "resttestkey" {
		t.Errorf("Expected APIKey to be 'resttestkey' but got '%s'", got)
	}
	if p.Provider.APIPassword != "" {
		t.Errorf("Expected APIPassword to be empty, got '%s'", p.Provider.APIPassword)
	}
	if p.Provider.CustomerNumber != "" {
		t.Errorf("Expected CustomerNumber to be empty, got '%s'", p.Provider.CustomerNumber)
	}
}

// TestUnmarshalCaddyFileInvalidCombinations verifies that a mix of legacy
// fields that is missing exactly one of customer_number/api_password is
// rejected, instead of silently falling back to REST-API mode (which
// supports far fewer record types) due to a typo or a forgotten field.
func TestUnmarshalCaddyFileInvalidCombinations(t *testing.T) {
	tests := []struct {
		name string
		conf string
	}{
		{
			name: "missing api_key",
			conf: `netcup {
				customer_number 12345
				api_password pw123
			}`,
		},
		{
			name: "api_password without customer_number",
			conf: `netcup {
				api_key apikey
				api_password pw123
			}`,
		},
		{
			name: "customer_number without api_password",
			conf: `netcup {
				customer_number 12345
				api_key apikey
			}`,
		},
		{
			name: "empty block",
			conf: `netcup {
			}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dispenser := caddyfile.NewTestDispenser(tc.conf)
			p := Provider{&netcup.Provider{}}
			if err := p.UnmarshalCaddyfile(dispenser); err == nil {
				t.Errorf("UnmarshalCaddyfile(%q) should have failed but didn't", tc.conf)
			}
		})
	}
}
