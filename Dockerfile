FROM golang:1.14

WORKDIR /go/src/github.com/haylesnortal/rpt/rpt
COPY ./rpt .

RUN go get -d -v ./...
RUN go install -v ./...

WORKDIR /go/src/app
COPY ./app.go .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]