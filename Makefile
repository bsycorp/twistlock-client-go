all: compile

clean:
	rm -rf build

compile: clean
	go mod download
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./build/twistlock-controller-linux-x64 ./cmd/twistlock-controller
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o ./build/twistlock-controller-darwin-x64 ./cmd/twistlock-controller
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o ./build/twistlock-controller-win-x64.exe ./cmd/twistlock-controller
