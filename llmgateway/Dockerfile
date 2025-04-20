FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/user_center

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/ .
EXPOSE 8080
CMD ["/app/bin/user_center", "--proxyConf", "conf/proxy.yml", "--resourceConf", "conf/resource.yml"]