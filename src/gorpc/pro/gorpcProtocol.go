package pro

import (
	"net"
	"gorpc/utils"
	"net/rpc"
	"log"
	"net/http"
	"reflect"
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

func NewHTTPServer(service []interface{}){
	for _ ,s := range service {
		log.Println("register http service:",reflect.TypeOf(s).String())
		rpc.Register(s)
	}
	rpc.HandleHTTP()
	err := http.ListenAndServe(":1234",nil)
	utils.CheckErr(err)
}

func NewHTTPClient(host string) *rpc.Client{
	client,err := rpc.DialHTTP("tcp" , host)
	utils.CheckErr(err)
	return client
}