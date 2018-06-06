FROM golang:1.10.2

WORKDIR /go/src/github.com/wjglerum/kube-crd
COPY . .

RUN go install -v ./...

CMD ["kube-crd", "-kubeconf=/config"]
