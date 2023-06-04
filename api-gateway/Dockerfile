FROM golang:latest

RUN apt-get update


WORKDIR /dev-cloud-share


COPY . .

RUN go build -o api-gateway .

EXPOSE 8440
CMD ["./api-gateway"]


