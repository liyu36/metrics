FROM docker.io/library/golang:1.21.4 as builder
#RUN go env -w GO111MODULE="on"
#RUN go env -w GOPROXY="https://goproxy.cn,direct"
RUN mkdir /app
WORKDIR /app
ADD . /app
RUN go mod tidy
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o main main.go

# FROM scratch
FROM registry.cn-hangzhou.aliyuncs.com/liy36/alpine:3.15
WORKDIR /app
COPY --from=builder /app/metrics /app/main
CMD ["/app/metrics"]
