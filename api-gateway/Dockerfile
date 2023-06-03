FROM golang:latest

RUN apt-get update


WORKDIR /dev-cloud-share


COPY . .

RUN go build -o api-gateway .


CMD ["./api-gateway"]


