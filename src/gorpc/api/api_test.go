package api

import (
	"testing"
	"log"
	"gorpc/utils"
	"net/rpc"
)

func TestIp(t *testing.T) {
	t.Log(utils.Ip())
}

func TestGoRpc_RegisterServer(t *testing.T) {
	rpc := NewGoRpc("http://127.0.0.1:2379")
	tes := &Test{}
	//tes.Tostring(Request{"123"},&Response{"123"})
	rpc.RegisterServer(tes)
	w := make(chan int)
	<- w
}

func TestGoRpc_Call(t *testing.T) {
	goRpc := NewGoRpc("http://127.0.0.1:2379")
	resp := &Response{}
	f := Facade{"api.Test","Tostring",Request{"request test!!!"},resp}
	goRpc.Call(f)
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
