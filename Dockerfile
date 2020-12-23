FROM golang:latest
WORKDIR /go/src/github.com/gDreamcatcher/rpcxdemo
COPY . .
ENTRYPOINT ["go", "run", "main.go", "-method"]
CMD ["server"]
