FROM golang:1.14-alpine AS build

WORKDIR /go/src/github.com/HackDalton/homemade-http
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine

COPY --from=build /go/bin/homemade-http ./
COPY ./static ./static
COPY flag.txt flag.txt

CMD ["./homemade-http"]