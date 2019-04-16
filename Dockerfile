FROM golang:1.12 AS build-env
WORKDIR /src
RUN useradd -u 10001 controller

# Populate the module cache based on the go.{mod,sum} files.
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

# Build
RUN make

# run
FROM scratch
WORKDIR /
COPY --from=build-env /src/build/twistlock-controller-linux-x64 /twistlock-controller
COPY --from=build-env /etc/passwd /etc/passwd
USER controller

ENTRYPOINT ["/twistlock-controller"]
