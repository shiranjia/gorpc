package main

import (
	"time"
	"gorpc/api"
	"log"
	"gorpc/utils"
)

func main() {
	goRpc := api.NewGoRpc("http://192.168.146.147:2379")
	func(){
		resp := &api.Response{}
		f := api.Facade{"api.Test1","Tostring",api.Request{"request test!!!"},resp}
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

	time.Sleep(100 * time.Second)
}


