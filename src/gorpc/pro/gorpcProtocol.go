package pro

import (
	"net"
	"gorpc/utils"
	"net/rpc"
	"log"
)

func NewServer(service interface{}){
	rpc.Register(service)
	listener,err := net.Listen("tcp","127.0.0.1:7777")
	utils.CheckErr(err)
	listen(listener)
	log.Println("register service:",service)
}

func listen(l net.Listener){
	go func(){
		for  {
			conn,err := l.Accept()
			utils.CheckErr(err)
			rpc.ServeConn(conn)
		}
	}()
}

func NewClient(host string) *rpc.Client{
	client,err := rpc.Dial("tcp" , host)
	utils.CheckErr(err)
	return client
}