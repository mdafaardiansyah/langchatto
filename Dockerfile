FROM golang:1.22.2-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

COPY .env ./

RUN go build -o langchatto-app ./cmd/main.go

RUN chmod +x langchatto-app

EXPOSE 4000

EXPOSE 8080

CMD ["./langchatto-app"]