package tw

type IsInitialized struct {
	Initialized bool `json:"initialized"`
}

func (c *Client) IsInitialized() (bool, error) {
	req, err := c.newRequest("GET", "settings/initialized", nil)
	if err != nil {
		return false, err
	}
	var resp IsInitialized
	_, err = c.do(req, &resp)
	if err != nil {
		return false, err
	}
	return resp.Initialized, nil
}
