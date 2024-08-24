FROM golang:1.23.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the plugin with CGO enabled
#RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -buildmode=plugin -o ./plugin/hook_plugin ./hooks/hook_plugin.go
RUN CGO_ENABLED=1 GOOS=linux go build -buildmode=plugin -o ./plugin/hook_plugin ./hooks/hook_plugin.go

RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go
#RUN go build -o main main.go

#Final stage: create a minimal image to run the application
FROM golang:1.23.0-alpine AS final

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/plugin ./plugin
COPY --from=builder /app/.env .

RUN chmod +x ./plugin/hook_plugin

EXPOSE 8080

CMD ["./main"]