package api

import (
	"testing"
	"log"
	"gorpc/utils"
)

func TestIp(t *testing.T) {
	t.Log(utils.Ip())
}

func TestGoRpc_RegisterServer(t *testing.T) {
	rpc := NewGoRpc("127.0.0.1")
	tes := &Test{}
	tes.Tostring(Request{"123"},&Response{"123"})
	rpc.RegisterServer("",tes)
	w := make(chan int)
	<- w
}

type Test struct {

}

func (t *Test) Tostring(req Request,resp *Response)  error {
	log.Println(req.Body)
	resp.Body = "test"
	return nil
}
