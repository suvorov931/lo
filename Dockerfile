FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod  ./
RUN go mod download

COPY . .
RUN go build -o lo cmd/*.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=build /app/lo .

CMD ["./lo"]