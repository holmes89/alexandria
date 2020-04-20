FROM golang:1.14-buster AS deps
RUN apt-get update \
    && apt-get install -y git gcc musl-dev ca-certificates libssl1.1 wget \
    && wget https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-18.04.deb \
    && dpkg -i seabolt-1.7.4-Linux-ubuntu-18.04.deb \
    && rm seabolt-1.7.4-Linux-ubuntu-18.04.deb \
    && apt-get clean

WORKDIR "/alexandria"
ADD server/*.mod server/*.sum ./
RUN go mod download

FROM deps AS build-env
ADD server/cmd ./cmd
ADD server/internal ./internal
ADD server/migrations ./migrations
ENV PORT 8080
EXPOSE 8080
RUN GOOS=linux go build -ldflags "-w -X main.docker=true" -o server cmd/*.go

FROM scratch AS prod
ENV PORT 8080
EXPOSE 8080
COPY --from=build-env /alexandria/server /
COPY --from=build-env /alexandria/migrations /migrations
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/server"]
