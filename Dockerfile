# Base build image
FROM golang:latest AS build_base

# set working directory inside the container
WORKDIR /kubernite

# copy in go dependency management files
COPY go.mod go.sum ./

# populate the module cache based on the go.{mod,sum} files
RUN go mod download

# due to the layer caching system in Docker the go mod download
# command will only be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction this line)

# the kubernite binary is build in this image
FROM build_base AS kubernite_builder

# copy the source from the current directory to the container's working directory
COPY . .

# build the Go app
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o kubernite ./cmd/kubernite/main.go

# this last stage produces the final build image
# start from a fresh Alpine image to reduce the image size
# (i.e. not ship the Go compiler in our production artifacts)
FROM alpine AS kubernite

# add git
RUN apk add --no-cache git

# add the certificates for TLS
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Finally we copy the statically compiled Go binary.
COPY --from=kubernite_builder /kubernite/kubernite /kubernite

ENTRYPOINT ["/kubernite"]
