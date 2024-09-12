# Build
FROM dockerhub.timeweb.cloud/library/golang:1.23 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./

ADD . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /gitch ./cmd/gitch

# Tests
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy
FROM dockerhub.timeweb.cloud/library/debian:11-slim AS build-release-stage

WORKDIR /

COPY --from=build-stage /gitch /gitch
RUN apt-get update
RUN apt-get install -y ca-certificates

EXPOSE 8080

ENTRYPOINT [ "/gitch"]