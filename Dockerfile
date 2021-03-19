FROM golang:1.16.2 as build-env

WORKDIR /go/src/holden
ADD . /go/src/holden

RUN go get -d -v ./...

RUN go build -o /go/bin/holden

FROM gcr.io/distroless/base
COPY ./assets/ /assets
COPY --from=build-env /go/bin/holden /

CMD ["/holden", "/config.toml"]
LABEL org.opencontainers.image.source=https://github.com/endigma/holden
