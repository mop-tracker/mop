FROM golang:alpine

RUN apk update && apk add --no-cache git
RUN git clone https://github.com/mop-tracker/mop ./mop
RUN cd mop && \
    go build ./cmd/mop && \
    chmod a+x ./mop
WORKDIR /go/mop/
CMD ["./mop"]
