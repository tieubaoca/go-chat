FROM golang:alpine
WORKDIR /app
COPY . .
RUN go mod tidy
CMD ["go", "run", "main.go", "start"]