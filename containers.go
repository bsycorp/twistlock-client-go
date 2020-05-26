package tw

import (
	"time"
)

/*
ContainerVulnerabilityResponse is for getting the vulnerabilities identified on container images
*/
type ContainerVulnerabilityResponse struct {
	ID                string   `json:"_id"`
	Collections       []string `json:"collections"`
	FirewallProtected bool     `json:"firewallProtected"`
	Hostname          string   `json:"hostname"`
	Info              struct {
		AllCompliance struct {
		} `json:"allCompliance"`
		App                    string `json:"app"`
		ComplianceDistribution struct {
			Critical int `json:"critical"`
			High     int `json:"high"`
			Low      int `json:"low"`
			Medium   int `json:"medium"`
			Total    int `json:"total"`
		} `json:"complianceDistribution"`
		ComplianceRiskScore       int `json:"complianceRiskScore"`
		ComplianceVulnerabilities []struct {
			Cause          string        `json:"cause"`
			Cve            string        `json:"cve"`
			Cvss           int           `json:"cvss"`
			Description    string        `json:"description"`
			Discovered     time.Time     `json:"discovered"`
			Exploit        string        `json:"exploit"`
			ID             int           `json:"id"`
			LayerTime      int           `json:"layerTime"`
			Link           string        `json:"link"`
			PackageName    string        `json:"packageName"`
			PackageVersion string        `json:"packageVersion"`
			Published      int           `json:"published"`
			RiskFactors    interface{}   `json:"riskFactors"`
			Severity       string        `json:"severity"`
			Status         string        `json:"status"`
			Templates      []interface{} `json:"templates"`
			Text           string        `json:"text"`
			Title          string        `json:"title"`
			Twistlock      bool          `json:"twistlock"`
			Type           string        `json:"type"`
			VecStr         string        `json:"vecStr"`
		} `json:"complianceIssues"`
		ComplianceVulnerabilitiesCnt int    `json:"complianceVulnerabilitiesCnt"`
		ID                           string `json:"id"`
		Image                        string `json:"image"`
		ImageID                      string `json:"imageID"`
		ImageName                    string `json:"imageName"`
		Infra                        bool   `json:"infra"`
		InstalledProducts            struct {
			Docker string `json:"docker"`
		} `json:"installedProducts"`
		Labels    []string `json:"labels"`
		Name      string   `json:"name"`
		Namespace string   `json:"namespace"`
		Network   struct {
			Ports []interface{} `json:"ports"`
		} `json:"network"`
		NetworkSettings struct {
		} `json:"networkSettings"`
		Processes []struct {
			Name string `json:"name"`
		} `json:"processes"`
		ProfileID string `json:"profileID"`
	} `json:"info"`
	ScanTime time.Time `json:"scanTime"`
}

/*
GetContainerList will gather a list of containers, and their attached vulns
the params argument turns into queryparams with their values.

Even though params is a map[string]string, and some values can be integers, this seems
to not matter ( I'm guessing this is a freebie from the http protocol )
*/
func (c *Client) GetContainerList(params map[string]string) ([]ContainerVulnerabilityResponse, error) {
	req, err := c.newRequest("GET", "containers", nil)
	if err != nil {
		return nil, err
	}
	if params != nil {
		parms := req.URL.Query()
		for queryParam, queryValue := range params {
			parms.Add(queryParam, queryValue)
		}
		req.URL.RawQuery = parms.Encode()
	}
	var ContainerVulnerabilities []ContainerVulnerabilityResponse
	_, err = c.do(req, &ContainerVulnerabilities)
	return ContainerVulnerabilities, err
}
