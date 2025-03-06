FROM golang:1.23
WORKDIR /app
COPY . .
RUN go build -o agent cmd/agent/agent.go
CMD ["./agent"]