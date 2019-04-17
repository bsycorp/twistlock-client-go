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
	ticker := time.NewTicker(15 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				err := ApplyConfig()
				if err != nil {
					log.Println("error applying config: ", err)
				}
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func ApplyConfig() error {
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
	}
	// OK, we have an admin user and a valid license!
	return nil
}