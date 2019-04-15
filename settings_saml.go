package tw

type SAMLSettings struct {
	AppID     string `json:"appId"`
	AppSecret struct {
		Encrypted string `json:"encrypted"`
		Plain string `json:"plain"`
	} `json:"appSecret"`
	Audience  string `json:"audience"`
	Cert      string `json:"cert"`
	ConsoleURL string `json:"consoleURL"`
	Enabled   bool   `json:"enabled"`
	Issuer    string `json:"issuer"`
	TenantID  string `json:"tenantId"`
	Type      string `json:"type"`
	URL       string `json:"url"`
}

func (c *Client) GetSAMLSettings() (*SAMLSettings, error) {
	req, err := c.newRequest("GET", "settings/saml", nil)
	if err != nil {
		return nil, err
	}
	var resp SAMLSettings
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil

}

func (c *Client) SetSAMLSettings(settings *SAMLSettings) error {
	req, err := c.newRequest("POST", "settings/saml", settings)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}
