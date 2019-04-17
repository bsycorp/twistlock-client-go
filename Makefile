all: compile

clean:
	rm -rf build

compile: clean
	GO111MODULE=on CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./build/twistlock-controller-linux-x64 ./cmd/twistlock-controller
