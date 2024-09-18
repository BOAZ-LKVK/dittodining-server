FROM golang:1.22.4 as golang

WORKDIR /dittodining
COPY . .

RUN cp -r ./vendor $GOPATH; exit 0

WORKDIR /dittodining/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/dittodining


FROM debian:stable-20230703-slim

COPY --from=golang /go/bin /app
ENTRYPOINT ["app/dittodining"]
