FROM golang:1.12.5 AS builder

WORKDIR /go/src/github.com/antonyho/go-project-demo

COPY . .

RUN make install

#----------------------------

FROM golang:1.12.5-stretch

WORKDIR /root/

COPY --from=builder /go/bin/websrv /go/bin/

EXPOSE 80

ENTRYPOINT ["websrv"]