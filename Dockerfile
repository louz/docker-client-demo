# build stage
FROM golang:1.10.3-alpine3.7 AS build-env
ADD . /src/docker-client-demo
ENV GOPATH /:/src/docker-client-demo/vendor
WORKDIR /src/docker-client-demo
RUN go build -o app


# test stage
#FROM golang:1.8-alpine3.6
#WORKDIR /src/docker-client-demo
#RUN go test


# release stage
FROM alpine:3.7
WORKDIR /app
EXPOSE 8080
COPY --from=build-env /src/docker-client-demo/app /app/
CMD ["./app"]