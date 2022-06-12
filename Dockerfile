FROM golang:1.18.3-bullseye
RUN mkdir -p /app
WORKDIR /app
COPY . .
# RUN mkdir -p /builded
EXPOSE 8080
RUN CGO_ENABLED=0 GOOS=linux go build main.go
ENTRYPOINT ["./main"]