package api

import (
	"testing"
	"gorpc/utils"
	"gorpc/service"
)

func TestGoRpc_RegisterServer(t *testing.T) {
	rpc := NewGoRpc("http://192.168.146.147:2379")
	tes := &Test{}
	rpc.RegisterServer(service.Service{tes,utils.PROCOTOL_RPC})
	w := make(chan int)
	<- w
}

func TestGoRpc_Call(t *testing.T) {
	goRpc := NewGoRpc("http://192.168.146.147:2379")
	resp := &Response{}
	f := Facade{
		Service:"api.Test",
		Method:"Tostring",
		Args:Request{"asdasdttt"},
		Response:resp,
		Protocol:utils.PROCOTOL_RPC,
	}
	goRpc.Call(f)
	t.Log(resp.Body)
}

func TestGoRpc_RegisterHTTPServer(t *testing.T) {
	rpc := NewGoRpc("http://192.168.146.147:2379")
	tes := &Test{}
	tes1 := &Test1{}
	rpc.RegisterServer(service.Service{tes,utils.PROTOCOL_HTTP},service.Service{tes1,utils.PROTOCOL_HTTP})
	w := make(chan int)
	<- w
}

func TestGoRpc_CallHTTP(t *testing.T) {
	goRpc := NewGoRpc("http://192.168.146.147:2379")
	func(){
		f := Facade{
			Service:"api.Test1",
			Method:"Tostring",
			Args:Request{"asdasdttt"},
			Response:&Response{},
			Protocol:utils.PROTOCOL_HTTP,
		}
		goRpc.Call(f)
		t.Log(f.Response)
		f.Args = Request{"asfafe!!!"}
		err := goRpc.Call(f)
		utils.CheckErr("TestGoRpc_CallHTTP",err)
		t.Log(f.Response)

		f.Args = Request{"yyyyyyttttttttttt!!!"}
		goRpc.Call(f)
		utils.CheckErr("TestGoRpc_CallHTTP",err)
		t.Log(f.Response)
	}()
	//执行测试用例时不能阻塞测试用例进程，否则rpc句柄会一直阻塞
	//time.Sleep(100 * time.Second)
}


