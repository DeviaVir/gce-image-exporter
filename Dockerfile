FROM golang:1.14 as build
# Must be enabled for "go test -race ..."
ENV CGO_ENABLED 1
ENV GOPATH /go
COPY . /go/src/github.com/DeviaVir/gce-image-exporter/
WORKDIR /go/src/github.com/DeviaVir/gce-image-exporter/
# Get test dependencies & run tests.
RUN go mod tidy
RUN go test -race -v ./...

# Build a fully statically linked image.
ENV CGO_ENABLED 0
# Build and put the git commit hash into the binary.
RUN go get \
      -v \
      -ldflags "-X github.com/m-lab/go/prometheusx.GitShortCommit=$(git log -1 --format=%h)" \
      github.com/DeviaVir/gce-image-exporter


# Now copy the cbif binary into a minimal base image.
FROM alpine
# Add package, but prevent saving APKINDEX files with -no-cache.
RUN apk add --no-cache ca-certificates
COPY --from=build /go/bin/gce-image-exporter /
ENV PATH /:$PATH
ENTRYPOINT ["/gce-image-exporter"]
