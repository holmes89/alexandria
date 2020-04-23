FROM golang:1.14-alpine AS base
RUN apk update \
    && apk upgrade \
    && apk add git gcc musl-dev ca-certificates openssl-libs-static cmake make g++ openssl-dev \
    && rm -rf /var/cache/apk/*

# Install seabolt
RUN git clone -b 1.7 https://github.com/neo4j-drivers/seabolt.git /seabolt
WORKDIR /seabolt/build
RUN cmake -D CMAKE_BUILD_TYPE=Release -D CMAKE_INSTALL_LIBDIR=lib .. && cmake --build . --target install

FROM base as deps
WORKDIR "/alexandria"
ADD server/*.mod server/*.sum ./
RUN go mod download

FROM deps AS build-env
ADD server/cmd ./cmd
ADD server/internal ./internal
ADD server/migrations ./migrations
ENV PORT 8080
EXPOSE 8080
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-w -X main.docker=true" --tags "seabolt_static" -o server cmd/*.go
CMD ["./server"]

FROM alpine AS prod
RUN apk update \
    && apk upgrade \
    && apk add ca-certificates openssl-libs-static openssl-dev \
    && rm -rf /var/cache/apk/*

WORKDIR /
ENV PORT 8080
EXPOSE 8080
COPY --from=build-env /alexandria/server /
COPY --from=build-env /alexandria/migrations /migrations
CMD ["/server"]
