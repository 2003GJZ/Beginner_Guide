FROM golang:1.23-alpine

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main .

EXPOSE 8080
CMD ["./main"]

