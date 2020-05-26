package tw

import "time"

type Impacts struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Total    int `json:"total"`
}

type HighestRiskFactors struct {
	Network             bool `json:"network,omitempty"`
	RootPrivilege       bool `json:"rootPrivilege,omitempty"`
	NoSecurityProfile   bool `json:"noSecurityProfile,omitempty"`
	PrivilegedContainer bool `json:"privilegedContainer,omitempty"`
}

type Vulnerability struct {
	Cve                    string             `json:"cve,omitempty"`
	Description            string             `json:"description,omitempty"`
	HighestRiskFactors     HighestRiskFactors `json:"highestRiskFactors,omitempty"`
	ImpactedPackages       []string           `json:"impactedPkgs,omitempty"`
	ImpactedResourcesCount int                `json:"impactedResourcesCnt,omitempty"`
	Link                   string             `json:"link,omitempty"`
	RiskFactors            struct {
		AttackComplexityLow struct {
		} `json:"Attack complexity: low,omitempty"`
		AttackVectorNetwork struct {
		} `json:"Attack vector: network,omitempty"`
		MediumSeverity struct {
		} `json:"Medium severity,omitempty"`
		CriticalSeverity struct {
		} `json:"Critical severity,omitempty"`
		DoS struct {
		} `json:"DoS,omitempty"`
		HasFix struct {
		} `json:"Has fix,omitempty"`
		RecentVulnerability struct {
		} `json:"Recent vulnerability,omitempty"`
		RemoteExecution struct {
		} `json:"Remote execution,omitempty"`
	} `json:"riskFactors,omitempty"`
	RiskScore int    `json:"riskScore,omitempty"`
	Status    string `json:"status"`
}

type VulnerabilityResponse struct {
	ID       string    `json:"_id"`
	Modified time.Time `json:"modified"`
	Images   struct {
		Impacted        Impacts         `json:"impacted"`
		Cves            Impacts         `json:"cves"`
		Vulnerabilities []Vulnerability `json:"vulnerabilities,omitempty"`
	} `json:"images,omitempty"`
	Hosts struct {
		Impacted        Impacts         `json:"impacted"`
		Cves            Impacts         `json:"cves"`
		Vulnerabilities []Vulnerability `json:"vulnerabilities,omitempty"`
	} `json:"hosts,omitempty"`
	Functions struct {
		Impacted        Impacts         `json:"impacted"`
		Cves            Impacts         `json:"cves"`
		Vulnerabilities []Vulnerability `json:"vulnerabilities,omitempty"`
	} `json:"functions,omitempty"`
	Containers struct {
		Impacted        Impacts         `json:"impacted"`
		Cves            Impacts         `json:"cves"`
		Vulnerabilities []Vulnerability `json:"vulnerabilities,omitempty"`
	} `json:"containers,omitempty"`
}

type RiskTree struct {
	Container string `json:"container,omitempty"`
	Factors   struct {
		Network           bool `json:"network"`
		NoSecurityProfile bool `json:"noSecurityProfile"`
		RootPrivilege     bool `json:"rootPrivilege"`
	} `json:"factors"`
	Host      string `json:"host,omitempty"`
	Image     string `json:"image,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type RiskTrees []RiskTree

type ImpactedResources struct {
	ID        string               `json:"_id"`
	RiskTrees map[string]RiskTrees `json:"riskTree"`
}

type VulnerabilityResponses []VulnerabilityResponse

func (c *Client) GetStatsVulnerabilities() (VulnerabilityResponses, error) {
	req, err := c.newRequest("GET", "stats/vulnerabilities", nil)
	if err != nil {
		return nil, err
	}
	var resp VulnerabilityResponses
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetStatsVulnerabilitiesImpactedResources(cve string) (ImpactedResources, error) {
	req, err := c.newRequest("GET", "stats/vulnerabilities/impacted-resources", nil)
	var resp ImpactedResources
	if err != nil {
		return resp, err
	}
	params := req.URL.Query()
	params.Add("cve", cve)
	req.URL.RawQuery = params.Encode()
	_, err = c.do(req, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
