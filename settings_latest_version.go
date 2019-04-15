package tw

type LatestVersion struct {
	LatestVersion string `json:"latestVersion"`
}

func (c *Client) GetLatestVersion() (string, error) {
	req, err := c.newRequest("GET", "settings/latest-version", nil)
	if err != nil {
		return "", err
	}
	var resp LatestVersion
	_, err = c.do(req, &resp)
	if err != nil {
		return "", err
	}
	return resp.LatestVersion, nil
}
