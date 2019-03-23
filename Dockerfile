FROM golang:latest as builder

ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src/github.com/irotoris/jobkickqd
COPY . .
RUN make build

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/irotoris/build/jobkickqd /usr/local/bin/jobkickqd
