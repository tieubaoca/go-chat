FROM golang:alpine
COPY . .
RUN go get ./...
RUN go run main.go initDb
CMD ["go", "run", "main.go", "start"]