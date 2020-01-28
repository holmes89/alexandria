FROM golang:alpine
WORKDIR "/covergen"
ADD cover-gen/*.go cover-gen/*.mod cover-gen/*.sum ./
RUN apk update \
    && apk upgrade \
    && apk add git gcc musl-dev ca-certificates ghostscript jpegoptim \
    && rm -rf /var/cache/apk/*
RUN go get ./...
ENV PORT 8080
EXPOSE 8080
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -X main.docker=true" -o server .
CMD ["/covergen/server"]
