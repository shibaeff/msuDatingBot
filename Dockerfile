FROM golang:latest

LABEL version="1.0"

RUN mkdir /go/src/echoBot
WORKDIR /go/src/echoBot
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep

#ADD ./main.go /go/src/app
#COPY ./Gopkg.toml /go/src/app

RUN dep ensure
# RUN go test -v
RUN go build ./cmd/main.go
EXPOSE 8080:8080
CMD ["./echoBot"]
