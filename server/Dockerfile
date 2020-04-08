FROM golang:1.13-alpine AS deps
RUN apk update \
    && apk upgrade \
    && apk add git gcc musl-dev ca-certificates \
    && rm -rf /var/cache/apk/*
WORKDIR "/alexandria"
ADD server/*.mod server/*.sum ./
RUN go get ./...

FROM deps AS build-env
ADD server/cmd ./cmd
ADD server/internal ./internal
ADD server/migrations ./migrations
ENV PORT 8080
EXPOSE 8080
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -X main.docker=true" -o server cmd/*.go

FROM scratch AS prod
ENV PORT 8080
EXPOSE 8080
COPY --from=build-env /alexandria/server /
COPY --from=build-env /alexandria/migrations /migrations
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/server"]
