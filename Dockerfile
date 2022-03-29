FROM golang:1.18.0-alpine3.15 as builder

RUN apk --no-cache add git

ENV CGO_ENABLED=0

WORKDIR /go/src/
ADD go.mod go.sum /go/src/
RUN go mod download

ADD main.go /go/src/
RUN go build -o /subfinder -ldflags="-s -w" .


FROM alpine:3.15
RUN apk add --no-cache ca-certificates
# in case we want to save about 5MB...
# FROM scratch

COPY --from=builder /subfinder /subfinder

ENTRYPOINT ["/subfinder"]
