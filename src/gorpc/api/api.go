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
)

type goRpc struct {
	serversCache map[string][]string `服务列表`
	register.Register		`注册中心`
	lock chan int			`更新服务缓存锁`
}

/**
构造服务api
 */
func NewGoRpc(host string) *goRpc{
	g := &goRpc{
		serversCache : make(map[string][]string),
		lock : make(chan int,1),
		Register : register.CreateEtcdRegister(host),
	}
	return  g
}

/**
注册rpc协议服务
 */
func (r *goRpc) RegisterRPCServer(service ...interface{})  {
	r.registerService(service)
	pro.NewRPCServer(service)
	log.Println("register over")
	r.TimeTicker()//打开心跳
}

/**
远程调用  go rpc协议
 */
func (r *goRpc) CallRPC(s Facade) error  {
	s.Service = "*"+s.Service
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
func (r *goRpc) RegisterHTTPServer(service ...interface{})  {
	r.registerService(service)
	go pro.NewHTTPServer(service)//注册服务
	log.Println("register over")
	r.TimeTicker()//打开心跳
}

/**
远程调用 http协议
 */
func (r *goRpc) CallHTTP(s Facade) error  {
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
func (r *goRpc) RegisterJsonServer(service ...interface{}){

}

/**
远程调用 json协议
 */
func (r *goRpc) CallJSON(s Facade) error    {

	return nil
}

/**
向注册中心注册服务
 */
func(r *goRpc) registerService(service []interface{}){
	for _,ser := range service{
		t :=reflect.TypeOf(ser)
		serviceName := t.String()
		log.Println(serviceName)
		r.Register.Set(serviceName + utils.Separator + utils.Ip() +":1234" , "")
	}
}

/**
从已注册服务列表中取出一个
 */
func (r *goRpc) getHost(s Facade) (string,error)  {
	hosts := r.serversCache[s.Service]
	if hosts == nil || len(hosts)==0{
		nodes,err := r.Register.GetChildren(s.Service)
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
	r.Subscribe(utils.Key2path(s.Service) , make(chan int), func(cl *client.Response) {
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

type Request struct {
	Body interface{}
}

type Response struct {
	Body interface{}
}


