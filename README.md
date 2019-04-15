# twistlock-client-go
Twistlock API Client for golang

## Testing

To run the tests, you need an access token and license. The integration
tests will spin up a temporary console container, initialize it and then
run the tests.

```
export TW_ACCESS_TOKEN=<your-access-token>
export TW_LICENSE=<your-license-key>
go test -v -tags integration ./test/integration/
```
