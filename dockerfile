FROM golang:1.18.3-alpine3.15

RUN mkdir /app

ADD /generate-csv.go /app
ADD /get-musics-ids.go /app

WORKDIR /app

RUN go env -w GO111MODULE=off

RUN GOOS=windows go build get-musics-ids.go
RUN GOOS=windows go build generate-csv.go
RUN go build get-musics-ids.go
RUN go build generate-csv.go