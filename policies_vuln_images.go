package tw

import "time"

type ContainerVulnerabilityPolicy struct {
	Rules []struct {
		Modified     time.Time `json:"modified"`
		Owner        string    `json:"owner"`
		Name         string    `json:"name"`
		PreviousName string    `json:"previousName"`
		Effect       string    `json:"effect"`
		Resources    struct {
			Containers []string `json:"containers"`
			Functions  []string `json:"functions"`
			Hosts      []string `json:"hosts"`
			Images     []string `json:"images"`
			Labels     []string `json:"labels"`
			Namespaces []string `json:"namespaces"`
			Services   []string `json:"services"`
		} `json:"resources"`
		Action    []string `json:"action,omitempty"`
		Condition struct {
			Readonly        bool          `json:"readonly"`
			Device          string        `json:"device"`
			Vulnerabilities []interface{} `json:"vulnerabilities"`
		} `json:"condition"`
		Group          []string `json:"group,omitempty"`
		AlertThreshold struct {
			Disabled bool `json:"disabled"`
			Value    int  `json:"value"`
		} `json:"alertThreshold"`
		BlockThreshold struct {
			Enabled bool `json:"enabled"`
			Value   int  `json:"value"`
		} `json:"blockThreshold"`
		CveRules []struct {
			Effect      string `json:"effect"`
			ID          string `json:"id"`
			Description string `json:"description"`
			Expiration  struct {
				Enabled bool      `json:"enabled"`
				Date    time.Time `json:"date"`
			} `json:"expiration"`
		} `json:"cveRules,omitempty"`
		GraceDays int `json:"graceDays"`
	} `json:"rules"`
	PolicyType string `json:"policyType"`
	ID         string `json:"_id"`
}

func (c *Client) GetContainerVulnerabilityPolicy() (*ContainerVulnerabilityPolicy, error) {
	req, err := c.newRequest("GET", "policies/vulnerability/images", nil)
	if err != nil {
		return nil, err
	}
	var resp ContainerVulnerabilityPolicy
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SetContainerVulnerabilityPolicy(policy *ContainerVulnerabilityPolicy) error {
	req, err := c.newRequest("GET", "policies/vulnerability/images", policy)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}
