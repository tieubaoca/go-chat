FROM golang:1.19
WORKDIR /docker/go/src/chat-server
COPY . .
RUN go get ./...
CMD ["go", "run", "main.go", "start"]