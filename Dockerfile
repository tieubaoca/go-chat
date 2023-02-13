FROM golang:alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go run main.go initDb
CMD ["go", "run", "main.go", "start"]