package api

import (
	"gorpc/pro"
	"gorpc/register"
	 "math/rand"
	"gorpc/utils"
	"reflect"
	"log"
	"strings"
	"net/rpc"
	"fmt"
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
	services := []interface{}{}
	for _,ser := range service {
		t := reflect.TypeOf(ser)
		serviceName := t.String()
		serviceName = strings.Replace(serviceName, "*", "", -1)
		log.Println(serviceName)
		r.Register.Set(serviceName + utils.Separator + utils.Ip() + ":1234", "1")
		services = append(services,ser)
	}
	pro.NewServer(services)
}

/**
远程调用  go rpc协议
 */
func (r *goRpc) CallRPC(s Facade) error  {
	hosts := r.serversCache[s.Service]
	var client *rpc.Client
	var method string
	if hosts == nil || len(hosts)==0{
		nodes,err := r.Register.GetChildren(s.Service)
		utils.CheckErr("api.Call",err)
		log.Println(nodes)
		defer func() error{
			if err:=recover();err!=nil{
				if fmt.Sprint(err) == utils.RANGE_ERROR{
					log.Println("call rpc error : no alive provider")
				}else {
					log.Println("call rpc error : ",err)
				}
			}
			return err
		}()
		client = pro.NewClient(nodes[0].Key)
		r.cacheServer(nodes,s)//缓存
	}else {
		host := hosts[rand.Int() % len(hosts)]
		client = pro.NewClient(host)
		log.Println("point cache host:", host,"method:",s.Method)
	}
	return client.Call(method,s.Args,s.Response)
}

/**
注册http协议服务
 */
func (r *goRpc) RegisterHTTPServer(service ...interface{})  {
	services := []interface{}{}
	for _,ser := range service{
		t :=reflect.TypeOf(ser)
		serviceName := t.String()
		serviceName = strings.Replace(serviceName,"*","",-1)
		log.Println(serviceName)
		r.Register.Set(serviceName + utils.Separator + utils.Ip() +":1234" , "")
		services = append(services,ser)
	}
	go pro.NewHTTPServer(services)//注册服务
	log.Println("register over")
	r.TimeTicker()//打开心跳
}

/**
远程调用 http协议
 */
func (r *goRpc) CallHTTP(s Facade) error  {
	hosts := r.serversCache[s.Service]
	var cli *rpc.Client
	if hosts == nil || len(hosts)==0{
		nodes,err := r.Register.GetChildren(s.Service)
		utils.CheckErr("api.CallHTTP",err)
		log.Println(nodes)
		if len(nodes)==0{
			e := errors.New("call rpc error : no alive provider")
			log.Println(e.Error())
			return e
		}
		defer func() error{
			if err:=recover();err!=nil{
					log.Println("call rpc error : ",err)
			}
			return err
		}()
		log.Println("call host :",nodes[0].Key)
		cli = pro.NewHTTPClient(nodes[0].Key)
		if hosts == nil {//初次调用 ，初始订阅服务变化
			go subscribe(s,r)
		}
		r.cacheServer(nodes,s)//缓存

	}else {
		host := hosts[rand.Int() % len(hosts)]
		cli = pro.NewHTTPClient(host)
		log.Println("point cache host:", host,"method:",s.Method)
	}
	class := strings.Split(s.Service,".")
	className := class[len(class)-1]
	defer cli.Close()
	//log.Println("cli.Call")
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
		log.Println("set")
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
		log.Println("delete")
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


