package api

import (
	"testing"
	"gorpc/utils"
)

func TescdtIp(t *testing.T) {
	t.Log(utils.Ip())
}

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
		Service:"api.Test1",
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
		resp := &Response{}
		f := Facade{
			Service:"api.Test1",
			Method:"Tostring",
			Args:Request{"asdasdttt"},
			Response:resp,
		}
		goRpc.CallHTTP(f)
		t.Log(resp.Body)
		f.Args = Request{"asfafe!!!"}
		err := goRpc.CallHTTP(f)
		utils.CheckErr("TestGoRpc_CallHTTP",err)
		t.Log(resp.Body)

		f.Args = Request{"yyyyyyttttttttttt!!!"}
		goRpc.CallHTTP(f)
		utils.CheckErr("TestGoRpc_CallHTTP",err)
		t.Log(resp.Body)
	}()
	//执行测试用例时不能阻塞测试用例进程，否则rpc句柄会一直阻塞
	//time.Sleep(100 * time.Second)
}


