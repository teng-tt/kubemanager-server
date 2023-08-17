FROM golang:apline as builder
WORKDIR /go/src/kubemanager-server
COPY . .

RUN go env -w GO111MODULE=on \
&& go env -w GOPROXY=https://goproxy.cn,direct \
&& go env -w CGO_ENABLED=0 \
&& go env \
&& go mod tidy \
&& go build -o server .

FROM apline:latest

LABEL MAINTAINER = "teng-tt"

WORKDIR /go/src/kubemanager
COPY --from=builder /go/src/kubemanager-server/config.yaml ./config.yaml
COPY --from=builder /go/src/kubemanager-server/.kube/config ./.kube/config
COPY --from=builder /go/src/kubemanager-server/server ./

EXPOSE 8080
ENTRYPOINT ./server