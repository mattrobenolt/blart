FROM golang:1.4-cross

RUN mkdir -p /go/src/blart
WORKDIR /go/src/blart

ENV CROSSPLATFORMS \
        linux/amd64 linux/386 linux/arm \
        darwin/amd64 darwin/386 \
        freebsd/amd64 freebsd/386 freebsd/arm \
        windows/amd64 windows/386

ENV GOARM 5

CMD set -x \
    && go-wrapper download \
    && for platform in $CROSSPLATFORMS; do \
            GOOS=${platform%/*} \
            GOARCH=${platform##*/} \
                go build -v -o bin/blart-${platform%/*}-${platform##*/}; \
    done
