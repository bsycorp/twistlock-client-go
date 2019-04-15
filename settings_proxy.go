package tw

type ProxyPassword struct {
	Encrypted string `json:"encrypted"`
	Plain     string `json:"plain"`
}

type ProxySettings struct {
	Ca        string        `json:"ca"`
	HttpProxy string        `json:"httpProxy"`
	NoProxy   string        `json:"noProxy"`
	Password  ProxyPassword `json:"password"`
	User      string        `json:"user"`
}

func (c *Client) GetProxy() (*ProxySettings, error) {
	req, err := c.newRequest("GET", "settings/proxy", nil)
	if err != nil {
		return nil, err
	}
	var resp ProxySettings
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SetProxy(r *ProxySettings) error {
	req, err := c.newRequest("POST", "settings/proxy", r)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}
