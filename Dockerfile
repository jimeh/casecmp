FROM golang:alpine

ADD . /go/src/github.com/jimeh/casecmp

RUN go install github.com/jimeh/casecmp

EXPOSE 8080
CMD ["/go/bin/casecmp", "--port", "8080"]



# FROM scratch
# ADD bin/casecmp_linux_amd64 /casecmp
# EXPOSE 8080
# VOLUME /data
# WORKDIR /
# CMD ["/casecmp", "--port", "8080"]
