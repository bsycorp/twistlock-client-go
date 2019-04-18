// +build integration

package integration

import (
	"github.com/bsycorp/twistlock-client-go"
	"reflect"
	"testing"
)

var bogusCert = `-----BEGIN CERTIFICATE-----
MIIC+zCCAeOgAwIBAgIJAIF9Dukg5bBtMA0GCSqGSIb3DQEBBQUAMBQxEjAQBgNV
BAMMCWZub3JkLmNvbTAeFw0xOTA0MTIwNjQ3MzFaFw0yOTA0MDkwNjQ3MzFaMBQx
EjAQBgNVBAMMCWZub3JkLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBALlu5UoI3I7xjhexlP8Klgo3ukBoyJgs+H+ph4BmcVJMyjwVnP5VgonzUyy7
GAPtN6YrtUTMM/a9Zp/BleI9TxBCfFk0dCqIyP8Gc+qUChe0bOyhgDNBoJ6Bnq6o
Vm2PFTi9oj4B4BtDzXG0Khg7cEhAPUbSJBxGTaoecHSkwBg/1AjgwzjUSkGE8UiC
rCN1jayImTtxebAhZdqWmxpklm8tFmZv8O8B/0BpAzmMVOU3eUVKu6Vc8fLzWya/
MM2Pvkkir0WPibo2ApyjNYIeYkBXy/hqAeTbDSsYDTCtgO0xwOUNpgcNhyium3Ll
J3b2UdhXBkc8Z1VTRDKWkaf8yvMCAwEAAaNQME4wHQYDVR0OBBYEFAg8LsnJ9pob
mWoeAkpZ056FhMyuMB8GA1UdIwQYMBaAFAg8LsnJ9pobmWoeAkpZ056FhMyuMAwG
A1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAL6cvDy63QkpduW3TYg1ihl
XRhd4xU9O7wvWJo+1rbVmyy7pqI4YuMDjWUBfmRBak4U4G2+80ZrTyCltxeQXE/T
3iA6Zz2atfuMJqKbwbvGwcJmAGwWvTJnC+iGfMqv9ug556wXjkbN9hJbGRf399jj
EO97TB0hUnIgvt4taU3obMw0xbHA+2AmAlZlcf7XEHbk3McwVJ/GfSgtXoWDEMcg
BXRw9rzpKaua+0oQ7/327EZWkflOYEpvBsInRqPspYB7wOeNO5NoLN9yG5wL7H5T
dDUpO0l6ViX9RzT1zPC610aKFLlVDe7smVGqvq0usfC2FUoex/eqdkwgclBZdW0=
-----END CERTIFICATE-----`

func TestSettingsSetAndGetProxy(t *testing.T) {
	settings := &tw.ProxySettings{
		Ca: bogusCert,
		HttpProxy: "http://proxies.int.v2.brkn.place:3128",
		NoProxy: "127.0.0.1",
		Password: tw.SecretValue{
			Plain: "swizzle",
		},
		User: "wibble",
	}

	// Set the proxy
	err := client.SetProxy(settings)
	if err != nil {
		t.Fatal(err)
	}
	// Read it back and verify it's the same-ish
	settings2, err := client.GetProxy()
	if err != nil {
		t.Fatal(err)
	}
	if settings.Ca != settings2.Ca {
		t.Fatal("ca does not match")
	}
	if settings.HttpProxy != settings2.HttpProxy {
		t.Fatal("HttpProxy does not match")
	}
	if settings.NoProxy != settings2.NoProxy {
		t.Fatal("NoProxy does not match")
	}
	// Password will come back, but encrypted
	if settings2.Password.Plain != "" || settings2.Password.Encrypted == "" {
		t.Fatal("unexpected proxy password value")
	}
	// Clear the proxy again
	err = client.SetProxy(&tw.ProxySettings{})
	if err != nil {
		t.Fatal(err)
	}
	// Check that it was cleared
	clearSettings, err := client.GetProxy()
	if !reflect.DeepEqual(clearSettings, &tw.ProxySettings{}) {
		t.Fatal("proxy settings not empty after clear")
	}
}
