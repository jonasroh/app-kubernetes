From golang:latest

workdir /app

COPY web-app .

run go build -o app

EXPOSE 8000

CMD ["./app"]