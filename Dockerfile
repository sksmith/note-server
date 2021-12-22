FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./note-server ./cmd

FROM scratch

WORKDIR /app

COPY --from=builder /app/note-server /usr/bin/

ENTRYPOINT ["note-server"]