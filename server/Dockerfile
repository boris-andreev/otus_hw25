#build stage
FROM golang:alpine AS builder

WORKDIR /app
COPY . .

WORKDIR /app/server
RUN go get -d -v ./...

WORKDIR /app
RUN CGO_ENABLED=0 go build -o ./bin/server server/cmd/.
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@latest

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/bin/server /server
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/postgres/migrations /postgres/migrations
COPY --from=builder /app/.env /.env
EXPOSE 5001

COPY run.sh /run.sh
RUN chmod +x /run.sh

ENTRYPOINT [ "/run.sh" ]

