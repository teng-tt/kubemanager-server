FROM golang:apline as builder
WORKDIR /go/src/kubmanager/server
COPY . .

RUN go env -w GO111MODULE=on \
&& go env -w GOPROXY=https://goproxy.cn,direct \
&& go env -w CGO_ENABLED=0 \
&& go env \
&& go mod tidy \
&& go build -o server .

FROM apline:latest

LABEL MAINTAINER = "teng-tt"

WORKDIR /go/src/kubmanager/server
COPY --from=builder /go/src/kubmanaer/server/config.yaml ./config.yaml
COPY --from=builder /go/src/kubmanaer/server/.kube/config ./.kube/config
COPY --from=builder /go/src/kubmanaer/server/server ./

EXPOSE 8080
ENTRYPOINT ./server