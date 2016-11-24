package api

import (
	"gorpc/pro"
	"gorpc/register"
	 "math/rand"
	"gorpc/utils"
	"reflect"
	"log"
	"strings"
	"errors"
	"github.com/coreos/etcd/client"
	"gorpc/service"
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
	rpcService := []interface{}{}
	httpService := []interface{}{}
	jsonService := []interface{}{}
	json2rpcService := []interface{}{}
	json2rpcHttpService := []interface{}{}
	for _,s := range service{
		switch s.Protocol {
		case utils.PROTOCOL_RPC		:	rpcService = append(rpcService,s.Servic)
		case utils.PROTOCOL_HTTP	:	httpService = append(httpService,s.Servic)
		case utils.PROTOCOL_JSON	:	jsonService = append(jsonService,s.Servic)
		case utils.PROTOCOL_JSON2RPC	:	json2rpcService = append(json2rpcService,s.Servic)
		case utils.PROTOCOL_JSON2RPCHTTP:	json2rpcHttpService = append(json2rpcHttpService,s.Servic)
		default:rpcService = append(rpcService,s.Servic)
		}
	}
	if len(rpcService) != 0 {r.registerRPCServer(rpcService,utils.PROTOCOL_RPC	)}
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
	host,_ := r.getHost(s)
	client := pro.NewRPCClient(host)
	class := strings.Split(s.Service,".")
	className := class[len(class)-1]
	defer client.Close()
	//log.Println("cli.Call")
	return client.Call(className + "." + s.Method,s.Args,s.Response)
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
	s.Service = "*"+s.Service
	host,_ := r.getHost(s)
	cli := pro.NewHTTPClient(host)
	class := strings.Split(s.Service,".")
	className := class[len(class)-1]
	defer cli.Close()
	return cli.Call(className + "." + s.Method,s.Args,s.Response)
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
	s.Service = "*"+s.Service
	host,_ := r.getHost(s)
	cli := pro.NewJSONClient(host)
	class := strings.Split(s.Service,".")
	className := class[len(class)-1]
	defer cli.Close()
	return cli.Call(className + "." + s.Method,s.Args,s.Response)
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
	s.Service = "*"+s.Service
	host,_ := r.getHost(s)
	cli := pro.NewJSON2Client(host)
	class := strings.Split(s.Service,".")
	className := class[len(class)-1]
	defer cli.Close()
	return cli.Call(className + "." + s.Method,s.Args,s.Response)
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
	s.Service = "*"+s.Service
	host,_ := r.getHost(s)
	cli := pro.NewHttpJson2rpcClient(host)
	class := strings.Split(s.Service,".")
	className := class[len(class)-1]
	defer cli.Close()
	return cli.Call(className + "." + s.Method,s.Args,s.Response)
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
func (r *goRpc) getHost(s Facade) (string,error)  {
	hosts := r.serversCache[s.Service]
	if hosts == nil || len(hosts)==0{
		nodes,err := r.register.GetChildren(s.Service)
		utils.CheckErr("api.CallHTTP",err)
		log.Println(nodes)
		if len(nodes)==0{
			e := errors.New("call rpc error : no alive provider:" + s.Service)
			log.Println(e.Error())
			return "",e
		}
		log.Println("call host :",nodes[0].Key)
		if hosts == nil {//初次调用 ，初始订阅服务变化
			go subscribe(s,r)
		}
		r.cacheServer(nodes,s)//缓存
		return nodes[0].Key,nil
	}else {
		host := hosts[rand.Int() % len(hosts)]
		log.Println("point cache host:", host,"method:",s.Method)
		return host,nil
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


