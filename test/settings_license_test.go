// +build integration

package integration

import (
	"github.com/zxcmx/twistlock-client-go"
	"testing"
)

func TestSettingsGetLicense(t *testing.T) {
	resp, err := client.GetLicense()
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.CustomerID) == 0 {
		t.Fatal("customer id not set")
	}
}

func TestSettingsSetInvalidLicense(t *testing.T) {
	err := client.SetLicense("THIS_IS_NOT_VALID")
	if err == nil {
		t.Fatal("should not be able to set invalid license")
	}
	e, ok := err.(*tw.ServerError)
	if !ok {
		t.Fatal("unexpected error type: ", err)
	}
	if e.StatusCode != 400 {
		t.Fatal("unexpected status code: ", e.StatusCode)
	}
	if e.Err != "invalid license" {
		t.Fatal("unexpected error message: ", e.Err)
	}
}
