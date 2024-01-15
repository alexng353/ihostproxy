ARG GOLANG_VERSION="1.21.6"

FROM golang:$GOLANG_VERSION-alpine as builder
WORKDIR /go/src/github.com/alexng353/ihostproxy

RUN apk --no-cache add tzdata gcc musl-dev

COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN CGO_ENABLED=1 GOOS=linux go build -a -o ./ihostproxy

FROM alpine:latest as runner
WORKDIR /app
COPY --from=builder /go/src/github.com/alexng353/ihostproxy/static ./static
COPY --from=builder /go/src/github.com/alexng353/ihostproxy/ihostproxy ./ihostproxy
ENTRYPOINT ["./ihostproxy"]
