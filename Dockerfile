FROM golang:1.14.1 as builder

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
WORKDIR /go/src/github.com/irotoris/jobkickqd
COPY . .
RUN make build

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/irotoris/jobkickqd/build/jobkickqd /usr/local/bin/jobkickqd
