FROM golang:alpine
WORKDIR /docker/go/src/chat-server
COPY . .
RUN go get ./...
CMD ["go", "run", "main.go", "start"]