FROM golang:1.11-alpine AS build

RUN apk add --no-cache git \
  && go get -u github.com/aws/aws-sdk-go-v2

ADD . /src
RUN cd /src && go build -o cloudfront-signer

FROM alpine
WORKDIR /app
COPY --from=build /src/cloudfront-signer /usr/local/bin/
ENTRYPOINT cloudfront-signer
