package tw

type IntelligenceSettings struct {
	WindowsFeedEnabled bool   `json:"windowsFeedEnabled"`
	Enabled            bool   `json:"enabled"`
	Address            string `json:"address"`
	Token              string `json:"token"`
}

func (c *Client) GetIntelligenceSettings() (*IntelligenceSettings, error) {
	req, err := c.newRequest("GET", "settings/intelligence", nil)
	if err != nil {
		return nil, err
	}
	var resp IntelligenceSettings
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil

}

func (c *Client) SetIntelligenceSettings(settings *IntelligenceSettings) error {
	req, err := c.newRequest("POST", "settings/intelligence", settings)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}
