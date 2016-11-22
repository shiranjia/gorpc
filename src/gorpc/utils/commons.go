package utils

import (
	"log"
	"net"
	"fmt"
	"os"
	"strings"
)

const(
	RANGE_ERROR	= "runtime error: index out of range"
	RootPath	= "/gorpc"	//根目录
	Separator	= "/"		//目录分隔符
	//// include get, set, delete, update, create, compareAndSwap,
	G		= "get"
	S		= "set"
	D		= "delete"
	U		= "update"
	C		= "create"
	E		= "expire"
	CompareAndSwap	="compareAndSwap"
)

func CheckErr(str string,e error) error {
	if e!= nil{
		log.Println(str,"Err:",e)
	}
	return e
}

func HandlerErr(e error , handler func(interface{}))  {
	if e!= nil{
		log.Println("Err:",e)
		if handler != nil{
			handler(e)
		}
	}
}

func Try(do func(),handler func(interface{}))  {
	defer func(){
		if err := recover();err != nil{
			handler(err)
		}
	}()
	do()
}

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

