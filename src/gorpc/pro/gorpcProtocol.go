package pro

import (
	"net"
	"gorpc/utils"
	"net/rpc"
	"log"
	"net/http"
	"reflect"
)

func NewRPCServer(service []interface{}){
	for _ ,s := range service {
		log.Println("register rpc service:",reflect.TypeOf(s).String())
		rpc.Register(s)
	}
	listener,err := net.Listen("tcp","127.0.0.1:1234")
	utils.CheckErr("gorpcProtocol.NewServer",err)
	go listen(listener)
}

func listen(l net.Listener){
	for {
		conn,err := l.Accept()
		utils.CheckErr("gorpcProtocol.listen",err)
		rpc.ServeConn(conn)
	}
}

func NewRPCClient(host string) *rpc.Client{
	client,err := rpc.Dial("tcp" , host)
	utils.CheckErr("gorpcProtocol.NewClient",err)
	return client
}

func NewHTTPServer(service []interface{}){
	for _ ,s := range service {
		log.Println("register http service:",reflect.TypeOf(s).String())
		rpc.Register(s)
	}
	rpc.HandleHTTP()
	err := http.ListenAndServe(":1234",nil)
	utils.CheckErr("gorpcProtocol.NewHTTPServer",err)
}

func NewHTTPClient(host string) *rpc.Client{
	client,err := rpc.DialHTTP("tcp" , host)
	utils.CheckErr("gorpcProtocol.NewHTTPClient",err)
	return client
}