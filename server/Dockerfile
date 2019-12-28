FROM golang:1.13-alpine AS build-env
WORKDIR "/alexandria"
ADD  server/*.go server/*.mod server/*.sum ./
RUN apk update \
    && apk upgrade \
    && apk add git gcc musl-dev ca-certificates \
    && rm -rf /var/cache/apk/*
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -X main.docker=true" -o server .

FROM scratch AS prod
ENV PORT 8080
COPY --from=build-env /alexandria/server /
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/server"]
