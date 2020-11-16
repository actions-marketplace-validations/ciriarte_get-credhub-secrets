FROM golang:1.15.5 as build

COPY . /go/src/get-credhub-secrets 
ENV CGO_ENABLED=0
ENV LD_FLAGS="-ldflags=-buildid= -trimpath"
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/get-credhub-secrets

RUN go fmt main.go
RUN go mod vendor
RUN go mod tidy
RUN go build

FROM ubuntu
COPY --from=build /go/src/get-credhub-secrets/get-credhub-secrets /usr/local/bin/get-credhub-secrets
ENTRYPOINT [ "get-credhub-secrets" ]