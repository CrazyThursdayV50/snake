FROM golang:1.24.1 AS builder

LABEL stage=gobuilder

RUN dpkg --add-architecture amd64 \
	&& apt-get update \
	&& apt-get install -y --no-install-recommends gcc-x86-64-linux-gnu libc6-dev-amd64-cross

ENV CC=x86_64-linux-gnu-gcc

ENV CGO_ENABLED=1
ENV GOPROXY=https://goproxy.cn,direct
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags '-linkmode "external" -extldflags "-static"' -o /app/main main.go

FROM alpine AS final
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/main /app/main
ENV CONFIG_NAME=config
ENTRYPOINT ["/bin/sh", "-c", "./main -d /app/config -c $CONFIG_NAME"]
