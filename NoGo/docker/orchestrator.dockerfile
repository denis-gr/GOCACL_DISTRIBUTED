FROM golang:1.23
WORKDIR /app
COPY . .
RUN go build -o orchestrator cmd/orchestrator/orchestrator.go
CMD ["./orchestrator"]