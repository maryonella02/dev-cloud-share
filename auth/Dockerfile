FROM golang:latest

RUN apt-get update


WORKDIR /dev-cloud-share


COPY . .

RUN go build -o auth .

EXPOSE 8443
CMD ["./auth"]
