FROM golang:1.19-alpine3.17
WORKDIR /app
COPY . .
RUN go mod tidy
CMD ["go", "run", "main.go", "start"]