# twistlock-client-go
Twistlock API Client for golang

## Summary

```go
import tw "github.com/zxcmx/twistlock-client-go"

...

client, err = tw.NewClient("http://localhost:8081/api/v1/")
if err != nil {
	log.Fatalln("error creating twistlock client: ", err)
}
err = client.Login("twadmin", "great_password")
if err != nil {
	log.Fatalln("failed to log into twistlock console: ", err)
}

...

```

## Testing

To run the live integration tests, you need an access token and license. The 
integration tests will spin up a temporary console container, initialize it and 
then run a few basic tests.

```
export TW_ACCESS_TOKEN=<your-access-token>
export TW_LICENSE=<your-license-key>
go test -v -tags integration ./test/integration/
```

---

# twistlock-controller

The twistlock-controller is a helper for automated twistlock provisioning and
configuration management.

It can be used as a "one shot" configurator during the twistlock provisioning 
process, or as a running controller to continuously manage twistlock console 
configuration.

An example configuration file can be found at: 
[cmd/twistlock-controller/sample.yml](cmd/twistlock-controller/sample.yml)

## Command-line and environment variables

```
Usage of ./twistlock-controller:
  -api string
    	API URL for Twistlock API (TWISTLOCK_API) (default "http://twistlock-console:8081/api/v1/")
  -config string
    	Path to configuration file (default "twistlock.yml")
  -license string
    	Twistlock console license (TWISTLOCK_LICENSE)
  -pass string
    	Twistlock console password (TWISTLOCK_PASS)
  -user string
    	Twistlock console username (TWISTLOCK_USER)
```

Environment variables take precedence over command-line arguments.

## Configuration

### Init block

If the SAML configuration includes block like:

```
init:
  enabled: true
```

The controller will check whether twistlock is initialized (via 
`/api/v1/settings/initialized`). If the console is not initialized yet, the 
controller will initialize it. This performs the following steps:

* Create a local account in the console for admin/superuser
  * Credentials for the account will be taken from the `-user` and `-pass`
    arguments or the corresponding `TWISTLOCK_USER` and `TWISTLOCK_PASS`
    environment variables.
* Install a license into the console. 
  * The license is taken from the `-license` argument or the 
    `TWISTLOCK_LICENSE` environment variable.

If the `init:` block is not specified and the controller is not initialized,
the controller will wait for manual initialization of the twistlock console
before applying any further configuration.

### SAML

Valid SAML provider types are ADFS, Azure, GSuite, Okta, Ping, Shibboleth.



