FROM golang:latest

RUN apt-get update

WORKDIR /dev-cloud-share

COPY . .

RUN go build -o resource-manager .

EXPOSE 8080
CMD ["./resource-manager"]
