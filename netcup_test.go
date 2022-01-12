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
