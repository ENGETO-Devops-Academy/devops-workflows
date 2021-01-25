FROM golang:1.14 as builder

WORKDIR /go/src/github.com/cermakm/example-apps/echoserver

COPY go.mod go.mod
COPY go.sum go.sum
RUN  go mod download

COPY main.go main.go
COPY middleware/ middleware/

RUN  go build -o ./bin/echoserver main.go

# ---

FROM registry.access.redhat.com/ubi8:latest

LABEL io.k8s.display-name="Echo Server"
LABEL name="Echo Server"
LABEL vendor="Engeto"
LABEL version="v0.0.0"
LABEL release="N/A"
LABEL summary="Kubernetes Example Application: Echo Server"
LABEL description="See summary"

COPY --from=builder /go/src/github.com/cermakm/example-apps/echoserver/bin/echoserver /opt/app-root/bin/echoserver

ENTRYPOINT [ "/opt/app-root/bin/echoserver" ]

