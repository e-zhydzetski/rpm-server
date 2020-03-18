# build stage
FROM golang:1.13.8-alpine3.11 AS builder
RUN apk update \
    && apk add --no-cache \
        git \
        ca-certificates \
        tzdata \
    && update-ca-certificates
RUN adduser -D -u 10001 appuser
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN cd cmd/rpm-server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o rpm-server

# final stage
FROM centos:7.7.1908
RUN yum install -y createrepo
EXPOSE 8080
ENV LISTEN_ADDR :8080
RUN mkdir -p /opt/repos/7/updates/x86_64
ENV PUSH_PATH /opt/repos/7/updates/x86_64
ENV REPOS_ROOT /opt/repos
ENTRYPOINT ["./rpm-server"]
COPY --from=builder /workspace/cmd/rpm-server/rpm-server ./rpm-server