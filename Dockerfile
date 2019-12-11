FROM golang:1.13.5 as builder
WORKDIR /go/src/istio.io/
RUN git clone --depth=1 https://github.com/istio/istio.git
WORKDIR /go/src/istio.io/istio/mixer/adapter/layer5
ADD . .
RUN cd cmd; CGO_ENABLED=0 GOOS=linux \
    go build -a -ldflags="-w -s" -installsuffix cgo -v -o /go/bin/layer5adapter .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /bin/
COPY --from=builder /go/bin/layer5adapter .
ENTRYPOINT [ "/bin/layer5adapter" ]
CMD [ "49494" ]