FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN set -xe \
    && CGO_ENABLED=0 GOOS=linux go build \
        -ldflags="\
          -X scn-tusd-server/services.VersionName=$(go list -m -f "{{ .Version }}" github.com/tus/tusd/v2) \
          -X 'scn-tusd-server/services.BuildDate=$(date --utc)'" \
        -o main main.go

FROM alpine:3.20.2 AS final

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]