package utils

import (
	"log"
	"net"
	"fmt"
	"os"
	"strings"
)

const(
	RootPath		= "/gorpc"	//根目录
	Separator		= "/"		//目录分隔符

	// include get, set, delete, update, create, compareAndSwap, expire  etcd可订阅事件
	G			= "get"
	S			= "set"
	D			= "delete"
	U			= "update"
	C			= "create"
	E			= "expire"
	CompareAndSwap		= "compareAndSwap"

	PROCOTOL_RPC		= "rpc"			//协议类型 默认rpc default
	PROTOCOL_HTTP		= "http"		//协议类型 tcp传输gob
	PROTOCOL_JSON		= "json"		//协议类型 tcp传输json
	PROTOCOL_JSON2RPC	= "json2rpc"		//协议类型 tcp传输json2rpc
	PROTOCOL_JSON2RPCHTTP	= "json2rpchttp"	//协议类型 http传输json2rpc
)

/**
处理错误
 */
func CheckErr(str string,e error) error {
	if e!= nil{
		log.Println(str,"Err:",e)
	}
	return e
}

/**
支持动态传的参数的错误处理
 */
func HandlerErr(e error , handler func(interface{}))  {
	if e!= nil{
		log.Println("Err:",e)
		if handler != nil{
			handler(e)
		}
	}
}

/**
模拟 try catch 语句块
 */
func Try(do func(),handler func(interface{}))  {
	defer func(){
		if err := recover();err != nil{
			handler(err)
		}
	}()
	do()
}

/**
获取本机ip
 */
func Ip() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
	return ""
}

func Path2key(path string) string{
	ps := strings.Split(path,Separator)
	l := len(ps)
	key := ps[0]
	if l > 0 {
		key = ps[l - 1]
	}
	return key
}

func Key2path(key string) string{
	return RootPath + Separator + key
}

func Host(protocol string) string {
	switch protocol {
	case PROCOTOL_RPC		:	return Ip() + ":1234"
	case PROTOCOL_HTTP		:	return Ip() + ":1235"
	case PROTOCOL_JSON		:	return Ip() + ":1236"
	case PROTOCOL_JSON2RPC		:	return Ip() + ":1237"
	case PROTOCOL_JSON2RPCHTTP	:	return Ip() + ":1238"
	default				:	return Ip() + ":1234"
	}
}

