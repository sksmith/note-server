FROM golang:alpine as builder

WORKDIR /app

COPY . .

ARG VER=NOT_SUPPLIED
ARG SHA1=NOT_SUPPLIED
ARG NOW=NOT_SUPPLIED

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X github.com/sksmith/note-server/config.AppVersion=$VER \
    -X github.com/sksmith/note-server/config.Sha1Version=$SHA1 \
    -X github.com/sksmith/note-server/config.BuildTime=$NOW" \
    -o ./note-server ./cmd

RUN apk add --update ca-certificates

FROM scratch

WORKDIR /app

COPY --from=builder /app/note-server /usr/bin/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["note-server"]
