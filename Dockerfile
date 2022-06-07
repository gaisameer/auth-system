FROM golang:1.18.3-alpine3.16
COPY . .
EXPOSE 8080
CMD ["go", "run", "main.go"]