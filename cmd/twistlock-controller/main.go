package main

import (
	"flag"
	tw "github.com/bsycorp/twistlock-client-go"
	"log"
	"os"
)

var apiUrl = flag.String("api", "http://twistlock-console:8081/api/v1/",
	"API URL for Twistlock API (TWISTLOCK_API)")
var username = flag.String("user", "", "Twistlock console username (TWISTLOCK_USER)")
var password = flag.String("pass", "", "Twistlock console password (TWISTLOCK_PASS)")
var license = flag.String("license", "", "Twistlock console license (TWISTLOCK_LICENSE)")
var config = flag.String("config", "twistlock.yml", "Path to configuration file")

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
	c, err := tw.NewClient(*apiUrl)
	if err != nil {
		log.Fatalln("error initializing twistlock client: ", err)
	}
	if *username != "" {
		err = c.Login(*username, *password)
		if err != nil {
			log.Fatalln("login error", err)
		}
	}
}
