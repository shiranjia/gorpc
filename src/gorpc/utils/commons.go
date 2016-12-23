/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package utils

import (
	"log"
	"net"
	"fmt"
	"os"
	"strings"
	"time"
)

const(
	RootPath		= "/gorpc"	//根目录
	Separator		= "/"		//目录分隔符
	Provider		= "provider"
	Consumer		= "consumer"
	// include get, set, delete, update, create, compareAndSwap, expire  etcd可订阅事件
	G			= "get"
	S			= "set"
	D			= "delete"
	U			= "update"
	C			= "create"
	E			= "expire"
	CompareAndSwap		= "compareAndSwap"

	PROTOCOL_RPC		= "rpc"			//协议类型 tcp传输gob 默认rpc default
	PROTOCOL_HTTP		= "http"		//协议类型 http传输gob
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

func ProviderPath(serviceName string) string  {
	return serviceName + Separator + "provider" + Separator
}

func Host(protocol string) string {
	switch protocol {
	case PROTOCOL_RPC		:	return Ip() + ":1234"
	case PROTOCOL_HTTP		:	return Ip() + ":1235"
	case PROTOCOL_JSON		:	return Ip() + ":1236"
	case PROTOCOL_JSON2RPC		:	return Ip() + ":1237"
	case PROTOCOL_JSON2RPCHTTP	:	return Ip() + ":1238"
	default				:	return Ip() + ":1234"
	}
}

func RoundSelector(count int) int  {
	a := time.Now().Second()
	//log.Println(a)
	a = a & (count - 1)
	return a
}

