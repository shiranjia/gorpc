package main

import (
	"gorpc/api"
	"gorpc/utils"
	"log"
	"time"
)

func main() {
	goRpc := api.NewGoRpc("http://192.168.146.147:2379")
	func(){
		resp := &api.Response{}
		f := api.Facade{
			Service:"api.Test1",
			Method:"Tostring",
			Args:api.Request{"asdasdttt"},
			Response:resp,
		}
		goRpc.CallHTTP(f)
		log.Println(resp.Body)
		f.Args = api.Request{"asfafe!!!"}
		err := goRpc.CallHTTP(f)
		utils.CheckErr("TestGoRpc_CallHTTP",err)
		log.Println(resp.Body)

		f.Args = api.Request{"yyyyyyttttttttttt!!!"}
		goRpc.CallHTTP(f)
		utils.CheckErr("TestGoRpc_CallHTTP",err)
		log.Println(resp.Body)
	}()
	//执行测试用例时不能阻塞测试用例进程，否则rpc句柄会一直阻塞
	time.Sleep(100 * time.Second)
}


