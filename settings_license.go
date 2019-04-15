package tw

import (
	"time"
)

type LicenseSettings struct {
	CustomerID      string `json:"customer_id"`
	CustomerEmail   string `json:"customer_email"`
	ContractID      string `json:"contract_id"`
	AccessToken     string `json:"access_token"`
	Type            string `json:"type"`
	Defenders       int    `json:"defenders"`
	DefenderDetails []struct {
		Category string `json:"category"`
		Count    int    `json:"count"`
	} `json:"defender_details"`
	IssueDate      time.Time `json:"issue_date"`
	ExpirationDate time.Time `json:"expiration_date"`
}

func (c *Client) GetLicense() (*LicenseSettings, error) {
	req, err := c.newRequest("GET", "settings/license", nil)
	if err != nil {
		return nil, err
	}
	var resp LicenseSettings
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SetLicense(licenseKey string) error {
	k := map[string]string{
		"key": licenseKey,
	}
	req, err := c.newRequest("POST", "settings/license", k)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

