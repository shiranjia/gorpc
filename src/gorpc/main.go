package main

import (
	"gorpc/api"
	"gorpc/utils"
	"time"
	"log"
)

func main() {

	rpc := api.NewGoRpc("http://192.168.146.147:2379")
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C{
		f := api.Facade{
			Service:"api.Test",
			Method:"Tostring",
			Args:api.Request{"ttt protocol rpc"},
			Response:&api.Response{},
			Protocol:utils.PROTOCOL_RPC,
			//Protocol:utils.PROTOCOL_JSON,
			//Protocol:utils.PROTOCOL_JSON2RPC,
			//Protocol:utils.PROTOCOL_JSON2RPCHTTP,
		}
		rpc.Call(f)
		log.Println(f.Response)
	}
}


