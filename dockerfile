FROM golang:latest
WORKDIR /app
COPY cmd cmd
COPY go.* ./
RUN go mod download
RUN go build -o main ./cmd
EXPOSE 3535
CMD ["./main"]
