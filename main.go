package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/gDreamcatcher/rpcxdemo/pb"
	"github.com/gin-gonic/gin"
	gateway "github.com/rpcxio/rpcx-gateway"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/codec"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var(
	addr = flag.String("addr", "localhost:8972", "server address")
	method = flag.String("method", "server", "server address")
)

type Arith int

func main() {
	flag.Parse()
	if *method == "server" {
		Server()
	}else if *method == "client" {
		Client()
	}else {
		HttpServer()
	}
}

func Server(){

	s := server.NewServer()
	s.Register(new(Arith), "")
	s.Serve("tcp", *addr)
}

func (t *Arith) Mul(ctx context.Context, args *pb.ProtoArgs, reply *pb.ProtoReply) error {
	reply.C = args.A * args.B
	log.Printf("call: %d * %d = %d\n", args.A, args.B, reply.C)
	return nil
}

func Client(){
	option := client.DefaultOption
	option.SerializeType = protocol.ProtoBuffer

	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	xclient := client.NewXClient("Arith", client.Failfast, client.RandomSelect, d, option)
	defer xclient.Close()

	args := &pb.ProtoArgs{A: 10, B: 20}
	reply := &pb.ProtoReply{}
	if err := xclient.Call(context.Background(), "Muls", args, reply); err != nil{
		fmt.Println(err)
	}
	log.Printf("%d * %d = %d", args.A, args.B, reply.C)
}

func HttpClient(){
	args := &pb.ProtoArgs{A: 10, B: 20}
	reply := &pb.ProtoReply{}

	cc := codec.MsgpackCodec{}
	data, _ := cc.Encode(args)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8972/", bytes.NewReader(data))
	if err != nil {
		log.Fatal("failed to create request: ", err)
	}
	h := req.Header
	h.Set(gateway.XMessageID, "10000")
	h.Set(gateway.XMessageType, "0")
	h.Set(gateway.XSerializeType, "3")
	h.Set(gateway.XServicePath, "Arith")
	h.Set(gateway.XServiceMethod, "Mul")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("failed to read response: ", err)
	}
	defer res.Body.Close()
	replyData, err := ioutil.ReadAll(res.Body)
	err = cc.Decode(replyData, reply)
	if err != nil {
		log.Fatal("failed to decode reply: ", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)
}

func HttpServer(){
	r := gin.Default()
	r.GET("/mul", CallMul)
	r.Run(":8080")
}

func CallMul(ctx *gin.Context){
	a, err := strconv.Atoi(ctx.Query("a"))
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(ctx.Query("b"))
	if err != nil {
		panic(err)
	}
	args := &pb.ProtoArgs{A: int32(a), B: int32(b)}
	reply := &pb.ProtoReply{}

	cc := codec.MsgpackCodec{}
	data, _ := cc.Encode(args)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8972/", bytes.NewReader(data))
	if err != nil {
		log.Fatal("failed to create request: ", err)
	}
	h := req.Header
	h.Set(gateway.XMessageID, "10000")
	h.Set(gateway.XMessageType, "0")
	h.Set(gateway.XSerializeType, "3")
	h.Set(gateway.XServicePath, "Arith")
	h.Set(gateway.XServiceMethod, "Mul")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("failed to read response: ", err)
	}
	defer res.Body.Close()
	replyData, err := ioutil.ReadAll(res.Body)
	err = cc.Decode(replyData, reply)
	if err != nil {
		log.Fatal("failed to decode reply: ", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "reply": reply.C})
}