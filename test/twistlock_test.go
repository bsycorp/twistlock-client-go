// +build integration

package integration

import (
	twistlock "github.com/zxcmx/twistlock-client-go"
	"log"
	"os"
)

var (
	client *twistlock.Client
)

func init() {
	base_url := os.Getenv("TWISTLOCK_URL")
	username := os.Getenv("TWISTLOCK_USER")
	password := os.Getenv("TWISTLOCK_PASS")

	if username == "" || password == "" || base_url == "" {
		log.Fatalln("Set TWISTLOCK_USER, TWISTLOCK_PASS, TWISTLOCK_URL for integration tests")
	}
	var err error
	client, err = twistlock.NewClient(base_url)
	if err != nil {
		log.Fatalln("error initializing twistlock client: ", err)
	}
	err = client.Login(username, password)
	if err != nil {
		log.Fatalln("login error: ", err)
	}
}
