// +build integration

package integration

import (
	"fmt"
	"github.com/ory/dockertest"
	twistlock "github.com/zxcmx/twistlock-client-go"
	"log"
	"os"
	"testing"
)

var client *twistlock.Client
var twAccessToken string
var twLicense string
var twUsername = "admin"
var twPassword = "admin123"

const ConsoleVersion = "19_03_317"

func TestMain(m *testing.M) {
	twAccessToken = os.Getenv("TW_ACCESS_TOKEN")
	twLicense = os.Getenv("TW_LICENSE")
	if twAccessToken == "" || twLicense == "" {
		log.Fatalln("must set TW_ACCESS_TOKEN and TW_LICENSE for integration tests")
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}
	// docker run -d -p 8081:8081 registry-auth.twistlock.com/tw_<access-token>/twistlock/console:console_19_03_311
	log.Println("starting twistlock console test container")
	consoleRepo := fmt.Sprintf("registry-auth.twistlock.com/tw_%s/twistlock/console", twAccessToken)
	consoleTag := "console_" + ConsoleVersion
	resource, err := pool.Run(consoleRepo, consoleTag, []string{})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}
	log.Println("waiting for twistlock console to be ready")

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		apiUrl := fmt.Sprintf("http://localhost:%s/api/v1/", resource.GetPort("8081/tcp"))
		client, err = twistlock.NewClient(apiUrl)
		return client.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	log.Println("twistlock console is ready")
	isInitialized, err := client.IsInitialized()
	if err != nil {
		log.Fatalf("failed to check initialization state of server: %s", err)
	}
	if isInitialized == true {
		log.Fatalln("expected server to be initialized, but it was not")
	}
	err = client.Signup(twUsername, twPassword)
	if err != nil {
		log.Fatalf("failed to initialize admin account: %s", err)
	}
	err = client.Login(twUsername, twPassword)
	if err != nil {
		log.Fatalf("failed to log in to twistlock console: %s", err)
	}
	log.Println("configured admin account")
	err = client.SetLicense(twLicense)
	if err != nil {
		log.Fatalf("failed to configure license: %s", err)
	}
	log.Println("configured license")

	// Now client should be pointing towards an initialized twistlock container.
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}
	log.Println("destroyed Twistlock console test container")

	os.Exit(code)
}
