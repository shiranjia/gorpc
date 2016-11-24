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
package api

import (
	"gorpc/pro"
	"gorpc/register"
	"gorpc/utils"
	"reflect"
	"log"
	"strings"
	"errors"
	"github.com/coreos/etcd/client"
	"gorpc/service"
	"net/rpc"
)

type goRpc struct {
	serversCache map[string][]string `服务列表`
	register register.Register		`注册中心`
	lock chan int			`更新服务缓存锁`
}

/**
构造服务api
 */
func NewGoRpc(host string) *goRpc{
	g := &goRpc{
		serversCache : make(map[string][]string),
		lock : make(chan int,1),
		register : register.CreateEtcdRegister(host),
	}
	return  g
}

/**
注册服务
 */
func  (r *goRpc) RegisterServer(service ...service.Service) {
	rpcService		:= []interface{}{}
	httpService		:= []interface{}{}
	jsonService		:= []interface{}{}
	json2rpcService 	:= []interface{}{}
	json2rpcHttpService 	:= []interface{}{}
	for _,s := range service{
		switch s.Protocol {
		case utils.PROTOCOL_RPC		:	rpcService = append(rpcService,s.Servic)
		case utils.PROTOCOL_HTTP	:	httpService = append(httpService,s.Servic)
		case utils.PROTOCOL_JSON	:	jsonService = append(jsonService,s.Servic)
		case utils.PROTOCOL_JSON2RPC	:	json2rpcService = append(json2rpcService,s.Servic)
		case utils.PROTOCOL_JSON2RPCHTTP:	json2rpcHttpService = append(json2rpcHttpService,s.Servic)
		default				:	rpcService = append(rpcService,s.Servic)
		}
	}
	if len(rpcService) != 0	{r.registerRPCServer(rpcService,utils.PROTOCOL_RPC	)}
	if len(httpService) != 0 {r.registerHTTPServer(httpService,utils.PROTOCOL_HTTP)}
	if len(jsonService) != 0 {r.registerJsonServer(jsonService,utils.PROTOCOL_JSON)}
	if len(json2rpcService) != 0 {r.registerJson2RpcServer(json2rpcService,utils.PROTOCOL_JSON2RPC)}
	if len(json2rpcHttpService) !=0 {r.registerJson2RpcHttpServer(json2rpcHttpService,utils.PROTOCOL_JSON2RPCHTTP)}
	log.Println("register over")
	r.register.TimeTicker()//打开心跳
}

/**
调用服务
 */
func (r *goRpc) Call(s Facade) error {
	switch s.Protocol {
	case utils.PROTOCOL_RPC		:	return r.callRPC(s)
	case utils.PROTOCOL_HTTP	:	return r.callHTTP(s)
	case utils.PROTOCOL_JSON	:	return r.callJson(s)
	case utils.PROTOCOL_JSON2RPC	:	return r.callJson2Rpc(s)
	case utils.PROTOCOL_JSON2RPCHTTP:	return r.callJson2RpcHttp(s)
	default				:	return r.callRPC(s)
	}
}

/**
注册rpc协议服务
 */
func (r *goRpc) registerRPCServer(service []interface{},protocol string)  {
	r.registerService(service,protocol)
	pro.NewRPCServer(service)

}

/**
远程调用  go rpc协议
 */
func (r *goRpc) callRPC(s Facade) error  {
	s.Service = "*" + s.Service
	host,index,err := r.getHost(s)
	if err != nil{
		return err
	}
	method := getMethodName(s)
	var cli *rpc.Client
	var e error
	utils.Try(func(){
		cli := pro.NewRPCClient(host)
		defer cli.Close()
		e = cli.Call(method,s.Args,s.Response)
	},func(err interface{}){
		log.Println("failover,fail host :",host)
		hosts := r.serversCache[s.Service]
		if len(hosts) == 0{
			log.Println("api.callRPC no alive service:",s.Service)
		}
		index = index +1
		if index < len(hosts){
			cli = pro.NewRPCClient(hosts[index])
		}else{
			cli = pro.NewRPCClient(hosts[0])
		}
		defer cli.Close()
		e = cli.Call(method,s.Args,s.Response)
	})

	return e
}

/**
注册http协议服务
 */
func (r *goRpc) registerHTTPServer(service []interface{},protocol string)  {
	r.registerService(service,protocol)
	go pro.NewHTTPServer(service)//注册服务
}

/**
远程调用 http协议
 */
func (r *goRpc) callHTTP(s Facade) error  {
	s.Service = "*" + s.Service
	host,index,err := r.getHost(s)
	if err != nil{
		return err
	}
	method := getMethodName(s)
	var cli *rpc.Client
	var e error
	utils.Try(func(){
		cli := pro.NewHTTPClient(host)
		defer cli.Close()
		e = cli.Call(method,s.Args,s.Response)
	},func(err interface{}){
		log.Println("failover,fail host :",host)
		hosts := r.serversCache[s.Service]
		if len(hosts) == 0{
			log.Println("api.callHTTP no alive service:",s.Service)
		}
		index = index +1
		if index < len(hosts){
			cli = pro.NewHTTPClient(hosts[index])
		}else{
			cli = pro.NewHTTPClient(hosts[0])
		}
		defer cli.Close()
		e = cli.Call(method,s.Args,s.Response)
	})

	return e
}

/**
注册服务 json协议
 */
func (r *goRpc) registerJsonServer(service []interface{},protocol string){
	r.registerService(service,protocol)
	pro.NewJSONServer(service)
}

/**
远程调用 json协议
 */
func (r *goRpc) callJson(s Facade) error    {
	s.Service = "*" + s.Service
	host,index,err := r.getHost(s)
	if err != nil{
		return err
	}
	method := getMethodName(s)
	var cli *rpc.Client
	var e error
	utils.Try(func(){
		cli := pro.NewJSONClient(host)
		defer cli.Close()
		e = cli.Call(method,s.Args,s.Response)
	},func(err interface{}){
		log.Println("failover,fail host :",host)
		hosts := r.serversCache[s.Service]
		if len(hosts) == 0{
			log.Println("api.callJson no alive service:",s.Service)
		}
		index = index +1
		if index < len(hosts){
			cli = pro.NewJSONClient(hosts[index])
		}else{
			cli = pro.NewJSONClient(hosts[0])
		}
		defer cli.Close()
		e = cli.Call(method,s.Args,s.Response)
	})

	return e
}

/**
注册服务 json2rpc 协议
 */
func (r *goRpc) registerJson2RpcServer(service []interface{},protocol string){
	r.registerService(service,protocol)
	pro.NewJSON2Server(service)
}

/**
远程调用 json2rpc 协议
 */
func (r *goRpc) callJson2Rpc(s Facade) error    {
	s.Service = "*" + s.Service
	host,index,err := r.getHost(s)
	if err != nil{
		return err
	}
	method := getMethodName(s)
	var cli *rpc.Client
	var e error
	utils.Try(func(){
		cli := pro.NewJSON2Client(host)
		defer cli.Close()
		e = cli.Call(method,s.Args,s.Response)
	},func(err interface{}){
		log.Println("failover,fail host :",host)
		hosts := r.serversCache[s.Service]
		if len(hosts) == 0{
			log.Println("api.callJson2Rpc no alive service:",s.Service)
		}
		index = index +1
		if index < len(hosts){
			cli = pro.NewJSON2Client(hosts[index])
		}else{
			cli = pro.NewJSON2Client(hosts[0])
		}
		defer cli.Close()
		e = cli.Call(method,s.Args,s.Response)
	})

	return e
}

/**
注册服务 json2rpc http 协议
 */
func (r *goRpc) registerJson2RpcHttpServer(service []interface{},protocol string){
	r.registerService(service,protocol)
	pro.NewHttpJson2rpcServer(service)
}

/**
远程调用 json2rpc http 协议
 */
func (r *goRpc) callJson2RpcHttp(s Facade) error    {
	s.Service = "*" + s.Service
	host,index,err := r.getHost(s)
	if err != nil{
		return err
	}
	method := getMethodName(s)
	var cli *rpc.Client
	var e error
	utils.Try(func(){
		cli := pro.NewHttpJson2rpcClient(host)
		defer cli.Close()
		e = cli.Call(method,s.Args,s.Response)
	},func(err interface{}){
		log.Println("failover,fail host :",host)
		hosts := r.serversCache[s.Service]
		if len(hosts) == 0{
			log.Println("api.callJson2RpcHttp no alive service:",s.Service)
		}
		index = index +1
		if index < len(hosts){
			cli = pro.NewHttpJson2rpcClient(hosts[index])
		}else{
			cli = pro.NewHttpJson2rpcClient(hosts[0])
		}
		defer cli.Close()
		e = cli.Call(method,s.Args,s.Response)
	})

	return e
}

func  getMethodName(s Facade) string{
	class := strings.Split(s.Service,".")
	return class[len(class)-1] + "." + s.Method
}

/**
向注册中心注册服务
 */
func(r *goRpc) registerService(service []interface{},protocol string){
	for _,ser := range service{
		t :=reflect.TypeOf(ser)
		serviceName := t.String()
		log.Println(serviceName)
		r.register.Set(serviceName + utils.Separator + utils.Host(protocol) , "")
	}
}

/**
从已注册服务列表中取出一个
 */
func (r *goRpc) getHost(s Facade) (string,int,error)  {

	hosts := r.serversCache[s.Service]
	if hosts == nil || len(hosts)==0{
		nodes,err := r.register.GetChildren(s.Service)
		utils.CheckErr("api.CallHTTP",err)
		log.Println(nodes)
		if len(nodes)==0{
			e := errors.New("call rpc error : no alive provider:" + s.Service)
			log.Println(e.Error())
			return "", 0 ,e
		}
		if hosts == nil {//初次调用 ，初始订阅服务变化
			go subscribe(s,r)
		}
		r.cacheServer(nodes,s)//缓存
		if len(nodes) == 1 {
			return nodes[0].Key,0,nil
		}else{
			index := utils.RoundSelector(len(nodes))
			return nodes[index].Key,index,nil
		}
	}else {
		host := ""
		index := 0
		if len(hosts) == 1 {
			host = hosts[0]
		}else {
			index = utils.RoundSelector(len(hosts))
			host = hosts[index]
		}
		log.Println("point cache host:", host,"method:",s.Method)
		return host,index,nil
	}
}

/**
缓存服务
 */
func  (r *goRpc)  cacheServer(nodes []register.Node,s Facade){
	//加入本地缓存
	log.Println("call method :",s.Method)
	hosts := []string{}
	for _ , h := range nodes {
		hosts = append(hosts,h.Key)
	}
	r.serversCache[s.Service] = append(r.serversCache[s.Service],hosts...)
	log.Println("serversCache :",r.serversCache)
}

/**
订阅服务注册中心
 */
func subscribe(s Facade,r *goRpc){
	r.register.Subscribe(utils.Key2path(s.Service) , make(chan int), func(cl *client.Response) {
		path := strings.Split(cl.Node.Key,utils.Separator)
		hostAndPort := path[len(path)-1]
		log.Println("收到事件：",cl.Action,cl.Node)
		switch cl.Action {

		case utils.S: r.updateServersCache(s.Service ,true ,hostAndPort)

		case utils.E: r.updateServersCache(s.Service ,false ,hostAndPort)

		}
	})
	log.Println("subscribe over")
}

/**
更新本地服务缓存
 */
func (r *goRpc) updateServersCache(serviceName string ,addOrDel bool,host string) {
	if addOrDel{
		log.Println(utils.S)
		add := true
		for _ , v := range r.serversCache[serviceName]{
			if v == host{
				add = false
				break//已缓存
			}
		}
		if add{
			r.lock <- 1
			r.serversCache[serviceName] = append(r.serversCache[serviceName],host)
			<- r.lock
		}
	}else {
		log.Println(utils.D)
		for i,v := range r.serversCache[serviceName] {
			if v == host{
				r.lock <- 1
				r.serversCache[serviceName] = append(r.serversCache[serviceName][:i],r.serversCache[serviceName][i+1:]...)
				<- r.lock
			}
		}
	}

	log.Println("serversCache changed:",serviceName,r.serversCache[serviceName])
}

/**
调用统计
 */
func (r *goRpc) invokStatistics( )  {

}

type Request struct {
	Body interface{}
}

type Response struct {
	Body interface{}
}


