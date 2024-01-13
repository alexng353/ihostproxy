ARG GOLANG_VERSION="1.21.6"

FROM golang:$GOLANG_VERSION-alpine as builder
RUN apk --no-cache add tzdata
WORKDIR /go/src/github.com/alexng353/ihostproxy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s' -o ./ihostproxy

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /go/src/github.com/alexng353/ihostproxy/ihostproxy /
ENTRYPOINT ["/ihostproxy"]
