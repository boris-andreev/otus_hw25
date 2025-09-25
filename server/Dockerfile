#build stage
FROM golang:alpine AS builder

WORKDIR /app
COPY . .

WORKDIR /app/server
RUN go get -d -v ./...

WORKDIR /app
RUN CGO_ENABLED=0 go build -o ./bin/server server/cmd/.

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/bin/server /server
EXPOSE 5001

CMD [ "/server" ]

