package main

import (
	"flag"
	tw "github.com/bsycorp/twistlock-client-go"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

var apiUrl = flag.String("api", "http://twistlock-console:8081/api/v1/",
	"API URL for Twistlock API (TWISTLOCK_API)")
var username = flag.String("user", "", "Twistlock console username (TWISTLOCK_USER)")
var password = flag.String("pass", "", "Twistlock console password (TWISTLOCK_PASS)")
var license = flag.String("license", "", "Twistlock console license (TWISTLOCK_LICENSE)")
var configFile = flag.String("config", "twistlock.yml", "Path to configuration file")

func main() {
	flag.Parse()
	if envApiUrl, ok := os.LookupEnv("TWISTLOCK_API"); ok {
		*apiUrl = envApiUrl
	}
	if envUser, ok := os.LookupEnv("TWISTLOCK_USER"); ok {
		*username = envUser
	}
	if envPass, ok := os.LookupEnv("TWISTLOCK_PASS"); ok {
		*password = envPass
	}
	if envLicense, ok := os.LookupEnv("TWISTLOCK_LICENSE"); ok {
		*license = envLicense
	}
	for {
		ApplyConfig()
		time.Sleep(15 * time.Minute)
	}
}

func ApplyConfig() {
	log.Println("applying config")
	err := Configure()
	if err != nil {
		log.Println("error applying config: ", err)
	} else {
		log.Println("configuration applied")
	}
}

func Configure() error {
	// Load config file
	config, err := LoadConfig(*configFile)
	if err != nil {
		return errors.Wrap(err, "failed to load config file")
	}
	// Light up the client and do a health check of the twistlock console
	c, err := tw.NewClient(*apiUrl)
	if err != nil {
		return errors.Wrap(err, "error initializing twistlock client")
	}
	if err = c.Ping(); err != nil {
		return errors.Wrap(err, "twistlock console api failed health check")
	}
	initDone, err := c.IsInitialized();
	if err != nil {
		return errors.Wrap(err, "could not check if console is initialized")
	}
	// Do we need to do initial setup of admin user?
	if !initDone {
		if !config.Init.Enabled {
			return errors.New("console admin user does not exist, but auto initialiation is disabled")
		}
		// OK, we need to create a new twistlock user
		log.Println("console not initialized; starting auto initialization")
		err := c.Signup(*username, *password)
		if err != nil {
			return errors.Wrap(err, "failed to create admin user")
		}
		log.Println("created admin user: ", *username)
	}
	// OK, we have an admin user, so we need to log in
	err = c.Login(*username, *password)
	if err != nil {
		return errors.Wrap(err, "login failed")
	}
	// Check the license
	twLicense, err := c.GetLicense()
	if err != nil {
		return errors.Wrap(err, "failed to check license")
	}
	if twLicense.CustomerID == "" {
		// We just check it's there, we don't check for expiry at this time.
		if !config.Init.Enabled {
			return errors.New("console license not installed, but auto initialization is disabled")
		}
		err = c.SetLicense(*license)
		if err != nil {
			return errors.Wrap(err, "failed to install license")
		}
		log.Println("installed license")

	}
	// At this point, we have an admin user, we're logging in and the license
	// is valid. Everything else is best-effort.
	if err = ConfigureSaml(c, &config.Saml); err != nil {
		log.Println("error configuring saml: ", err)
	}
	if err = ConfigureProxy(c, &config.Proxy); err != nil {
		log.Println("error configuring proxy: ", err)
	}

	return nil
}

func ConfigureSaml(client *tw.Client, saml *SamlConfig) error {
	if !saml.Enabled {
		log.Println("not configuring saml")
		return nil
	}
	samlTypes := []string{"adfs", "azure", "gsuite", "okta", "ping", "shibboleth"}
	if !ListContains(samlTypes, saml.Type) {
		return errors.Errorf("invalid saml.type: %s", saml.Type)
	}
	if saml.URL == "" {
		return errors.Errorf("saml.url required but not specified")
	}
	if saml.Cert == "" {
		return errors.New("saml.cert required but not specified")
	}
	if saml.Issuer == "" {
		return errors.New("saml.issuer required but not specified")
	}
	err := client.SetSAMLSettings(&tw.SAMLSettings{
		Enabled: true,
		Type: saml.Type,
		URL: saml.URL,
		Audience: saml.Audience,
		Cert: saml.Cert,
		ConsoleURL: saml.ConsoleURL,
		Issuer: saml.Issuer,
		TenantID: saml.TenantID,
		AppID: saml.AppID,
		AppSecret: tw.SecretValue{
			Plain: saml.AppSecret,
		},
	})
	return err
}

func ConfigureProxy(client *tw.Client, pc *ProxyConfig) error {
	if pc.HttpProxy == "" {
		return nil
	}
	err := client.SetProxy(&tw.ProxySettings{
		Ca: pc.CaCert,
		HttpProxy: pc.HttpProxy,
		NoProxy: pc.NoProxy,
		User: pc.ProxyUser,
		Password: tw.SecretValue{
			Plain: pc.ProxyPass,
			Encrypted: "",
		},
	})
	return err
}

func ListContains(haystack []string, needle string) (bool) {
	if haystack == nil {
		return false
	}
	for _, straw := range haystack {
		if straw == needle {
			return true
		}
	}
	return false
}
