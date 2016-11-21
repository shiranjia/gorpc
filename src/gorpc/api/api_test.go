package api

import (
	"testing"
	"log"
	"gorpc/utils"
	"net/rpc"
)

func TescdtIp(t *testing.T) {
	t.Log(utils.Ip())
}

func TestGoRpc_RegisterServer(t *testing.T) {
	rpc := NewGoRpc("http://192.168.146.147:2379")
	tes := &Test{}
	rpc.RegisterServer(tes)
	w := make(chan int)
	<- w
}

func TestGoRpc_Call(t *testing.T) {
	goRpc := NewGoRpc("http://192.168.146.147:2379")
	resp := &Response{}
	f := Facade{"main.Test","Tostring",Request{"request test!!!"},resp}
	goRpc.Call(f)
	t.Log(resp.Body)
}

func TestGoRpc_RegisterHTTPServer(t *testing.T) {
	rpc := NewGoRpc("http://192.168.146.147:2379")
	tes := &Test{}
	tes1 := &Test1{}
	rpc.RegisterHTTPServer(tes,tes1)
	w := make(chan int)
	<- w
}

func TestGoRpc_CallHTTP(t *testing.T) {
	resp := &Response{}
	f := Facade{"api.Test1","Tostring",Request{"request test!!!"},resp}

	goRpc := NewGoRpc("http://192.168.146.147:2379")
	goRpc.CallHTTP(f)
	t.Log(resp.Body)
	f.Args = Request{"asfafe!!!"}
	goRpc.CallHTTP(f)
	t.Log(resp.Body)
}

func TestGoRpc_Call2RPC(t *testing.T) {
	client,err := rpc.Dial("tcp" , "127.0.0.1:7777")
	resp := new(Response)
	res  := new(Request)
	res.Body = "resquest test"
	utils.CheckErr(err)
	client.Call("Test.Tostring",res,resp)
}

type Test struct {}
func (t *Test) Tostring(req Request,resp *Response)  error {
	log.Println(req.Body)
	resp.Body = "test"
	return nil
}
type Test1 struct {}
func (t *Test1) Tostring(req Request,resp *Response)  error {
	log.Println(req.Body)
	resp.Body = "test1"
	return nil
}
