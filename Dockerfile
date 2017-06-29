FROM golang:1.8-alpine as builder
ADD . /go/src/github.com/jimeh/casecmp
WORKDIR /go/src/github.com/jimeh/casecmp
RUN CGO_ENABLED=0 go build -a -o /casecmp \
    -ldflags "-X main.Version=$(cat VERSION)"

FROM scratch
COPY --from=builder /casecmp /
EXPOSE 8080
WORKDIR /
CMD ["/casecmp", "--port", "8080"]
