FROM golang:1.19-alpine3.17
WORKDIR /app
ENV PRODUCTION=false
ENV PORT=8888
COPY . .
RUN go mod tidy
CMD go run main.go start --security="$PRODUCTION" --port="$PORT"