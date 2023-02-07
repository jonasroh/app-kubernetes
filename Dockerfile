From golang:latest

workdir /app

COPY web-app ./

run go build -o web-app/main.go ./

EXPOSE 8000

CMD ["web-app/main.go"]