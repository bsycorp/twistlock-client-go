package main

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Init InitConfig `yaml:"init"`
	Saml SamlConfig `yaml:"saml"`
}

type InitConfig struct {
	Enabled bool `yaml:"enabled"`
}

type SamlConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Type       string `yaml:"type"`
	URL        string `yaml:"url"`
	Cert       string `yaml:"cert"`
	Issuer     string `yaml:"issuer"`
	Audience   string `yaml:"audience"`
	ConsoleURL string `yaml:"consoleURL"`
	TenantID   string `yaml:"tenantId"`
	AppID      string `yaml:"appId"`
	AppSecret  string `yaml:"appSecret"`
}
//var ValidSamlTypes = []string{"adfs", "azure", "gsuite", "okta", "ping", "shibboleth"}

func LoadConfig(path string) (*Config, error) {
	var c Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "read error: ")
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, errors.Wrap(err, "parse error: ")
	}
	return &c, nil
}
