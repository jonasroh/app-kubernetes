From golang:latest

workdir /app

COPY go.mod ./
COPY main.go ./

run go build -o main.go ./

EXPOSE 8080

CMD ["./main.go"]