FROM golang:latest

RUN apt-get update


WORKDIR /dev-cloud-share


COPY . .

RUN go build -o containerization-engine .

EXPOSE 8082
CMD ["./containerization-engine"]
