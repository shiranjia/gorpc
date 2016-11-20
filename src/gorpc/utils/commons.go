package utils

import (
	"log"
	"net"
	"fmt"
	"os"
)

func CheckErr(e error) error {
	if e!= nil{
		log.Println("Err:",e)
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
