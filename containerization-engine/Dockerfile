FROM golang:latest

RUN apt-get update


WORKDIR /dev-cloud-share


COPY . .

RUN go build -o containerization-engine .

CMD ["./containerization-engine"]
