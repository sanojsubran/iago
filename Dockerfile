# builder image
FROM golang:1.18.2-alpine3.16 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o iago .

#run image
FROM alpine:3.14
COPY --from=builder /build/iago .
ENTRYPOINT [ "./iago" ]
