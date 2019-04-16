package main

import (
	"log"
	"net/url"
	"strings"
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

var ValidSamlTypes = []string{"adfs", "azure", "gsuite", "okta", "ping", "shibboleth"}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func MustBeOneOf(s string, values []string, msg string) {
	for _, v := range values {
		if v == s {
			return
		}
	}
	log.Fatalf("%s: value should be one of: %v", msg, values)
}

func MustNotBeMissing(s, msg string) {
	if s == "" {
		log.Fatalln(msg)
	}
}

func MustBeUrl(u, msg string) {
	if u == "" {
		return
	}
	if !IsUrl(u) {
		log.Fatalln(msg)
	}
}

func MustBeValidSamlConfig(config *SamlConfig) {
	if config.Enabled != true {
		return
	}
	samlType := strings.ToLower(config.Type)
	MustNotBeMissing(samlType, "saml 'type' is required")
	MustBeOneOf(samlType, ValidSamlTypes, "invalid saml 'type'")
	MustNotBeMissing(config.URL, "saml 'url' is required")
	MustNotBeMissing(config.Cert, "saml 'cert' is required")
	MustNotBeMissing(config.Issuer, "saml 'issuer' is required")
	MustNotBeMissing(config.Audience, "saml 'audience' is required")
	MustBeUrl(config.URL, "saml 'url' is invalid")
	if config.ConsoleURL != "" {
		MustBeUrl(config.ConsoleURL, "saml 'consoleURL' is invalid")
	}
}

