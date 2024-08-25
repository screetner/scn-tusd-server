FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN set -xe \
	&& CGO_ENABLED=0 GOOS=linux go build \
        -ldflags="-X 'scn-tusd-server/services.BuildDate=$(date --utc)'" \
        -o main main.go

FROM golang:1.23.0-alpine AS final

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]