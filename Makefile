all: compile

clean:
	rm -rf build

compile: clean
	GOARCH=amd64 GOOS=linux go build -o ./build/twistlock-controller-linux-x64 ./cmd/twistlock-controller
	GOARCH=amd64 GOOS=darwin go build -o ./build/twistlock-controller-darwin-x64 ./cmd/twistlock-controller
	GOARCH=amd64 GOOS=windows go build -o ./build/twistlock-controller-win-x64.exe ./cmd/twistlock-controller
