# Typical usage:
#   docker build --progress=plain --pull --rm -f "Dockerfile" -t slot:latest "."
#   docker run -d -p 8080:8080 slot

##
## Build stage
##

# Use image with golang last version as builder.
FROM golang:1.22-bullseye AS build

# See https://stackoverflow.com/questions/64462922/docker-multi-stage-build-go-image-x509-certificate-signed-by-unknown-authorit
RUN apt-get update && apt-get install -y ca-certificates openssl
ARG cert_location=/usr/local/share/ca-certificates
# Get certificate from "github.com".
RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt
# Get certificate from "proxy.golang.org".
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt
# Update certificates.
RUN update-ca-certificates

# Make project root folder as current dir.
WORKDIR /go/src/github.com/slotopol/server
# Copy only go.mod and go.sum to prevent downloads all dependencies again on any code changes.
COPY go.mod go.sum ./
# Download all dependencies pointed at go.mod file.
RUN go mod download
# Copy all files and subfolders in current state as is.
COPY . .
COPY ./cmd ./cmd
COPY ./spi ./spi
COPY ./config ./config
COPY ./confdata /go/bin/config

# Compile project for Linux amd64.
RUN \
buildvers=$(git describe --tags) && \
buildtime=$(go run ./task/timenow.go) && \
go env -w GOOS=linux GOARCH=amd64 CGO_ENABLED=1 && \
go build -o /go/bin/slot_linux_x64 -v -tags="jsoniter prod full" -ldflags="-linkmode external -extldflags -static -X 'github.com/slotopol/server/config.BuildVers=$buildvers' -X 'github.com/slotopol/server/config.BuildTime=$buildtime'" ./

##
## Deploy stage
##

# Thin deploy image.
FROM scratch

# Copy compiled executable and packages to new image destination.
COPY --from=build /go/bin/slot* /go/bin/
# Copy configuration files.
COPY --from=build /go/bin/config /go/bin/config

# Open REST listen port.
EXPOSE 8080

# Run application with full path representation.
# Without shell to get signal for graceful shutdown.
ENTRYPOINT ["/go/bin/slot_linux_x64", "web"]
