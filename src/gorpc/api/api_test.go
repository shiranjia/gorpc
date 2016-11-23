package api

import (
	"testing"
	"gorpc/utils"
)

func TestGoRpc_RegisterServer(t *testing.T) {
	rpc := NewGoRpc("http://192.168.146.147:2379")
	tes := &Test{}
	rpc.RegisterRPCServer(tes)
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
	}
	goRpc.CallRPC(f)
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
	goRpc := NewGoRpc("http://192.168.146.147:2379")
	func(){
		f := Facade{
			Service:"api.Test1",
			Method:"Tostring",
			Args:Request{"asdasdttt"},
			Response:&Response{},
		}
		goRpc.CallHTTP(f)
		t.Log(f.Response)
		f.Args = Request{"asfafe!!!"}
		err := goRpc.CallHTTP(f)
		utils.CheckErr("TestGoRpc_CallHTTP",err)
		t.Log(f.Response)

		f.Args = Request{"yyyyyyttttttttttt!!!"}
		goRpc.CallHTTP(f)
		utils.CheckErr("TestGoRpc_CallHTTP",err)
		t.Log(f.Response)
	}()
	//执行测试用例时不能阻塞测试用例进程，否则rpc句柄会一直阻塞
	//time.Sleep(100 * time.Second)
}


