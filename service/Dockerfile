#This will be used as a builder for the later docker image,
#this way the later container will have a very small size as it will contain only
#the necessary executable to run the service.
FROM golang:1.16-alpine as builder

WORKDIR /app

# COPY go.mod, go.sum and download the dependencies
COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o /app/build /app/cmd/main.go

#This will be the container for running the service
FROM alpine:latest as service-container

WORKDIR /app

COPY --from=builder /app/build /app/

EXPOSE 8080

CMD [ "/app/main" ]