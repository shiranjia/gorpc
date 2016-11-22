package api

import (
	"gorpc/pro"
	"gorpc/server"
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
	point server.Provider
	register.Register
	lock chan int

}

func NewGoRpc(host string) *goRpc{
	g := &goRpc{
		serversCache : make(map[string][]string),
		lock : make(chan int,1),
		Register : register.CreateEtcdRegister(host),
	}
	return  g
}

func (r *goRpc) RegisterServer(service ...interface{})  {
	services := []interface{}{}
	for _,ser := range service {
		t := reflect.TypeOf(ser)
		serviceName := t.String()
		serviceName = strings.Replace(serviceName, "*", "", -1)
		log.Println(serviceName)
		r.Register.Set(serviceName + "/" + utils.Ip() + ":1234", "1")
		services = append(services,ser)
	}
	pro.NewServer(services)
}

func (r *goRpc) Call(s Facade) error  {
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
		hosts := []string{}
		for _ , h := range nodes {
			hosts = append(hosts,h.Key)
		}
		r.serversCache[s.Service] = append(r.serversCache[s.Service],hosts...)
	}else {
		host := hosts[rand.Int() % len(hosts)]
		client = pro.NewClient(host)
		log.Println("point cache host:", host,"method:",s.Method)
	}
	return client.Call(method,s.Args,s.Response)
}

func (r *goRpc) RegisterHTTPServer(service ...interface{})  {
	services := []interface{}{}
	for _,ser := range service{
		t :=reflect.TypeOf(ser)
		serviceName := t.String()
		serviceName = strings.Replace(serviceName,"*","",-1)
		log.Println(serviceName)
		r.Register.Set(serviceName + "/" + utils.Ip() +":1234" , "")
		services = append(services,ser)
	}
	go pro.NewHTTPServer(services)//注册服务
	log.Println("register over")
	r.TimeTicker()//打开心跳
}

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
		//加入本地缓存
		log.Println("call method :",s.Method)
		hosts := []string{}
		for _ , h := range nodes {
			hosts = append(hosts,h.Key)
		}
		r.serversCache[s.Service] = append(r.serversCache[s.Service],hosts...)
		log.Println("serversCache :",r.serversCache)
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

func subscribe(s Facade,r *goRpc){
	r.Subscribe(utils.Key2path(s.Service),make(chan int), func(cl *client.Response) {
		path := strings.Split(cl.Node.Key,"/")
		hostAndPort := path[len(path)-1]
		log.Println()
		switch cl.Action {

		case utils.S: r.updateServersCache(s.Service ,true ,hostAndPort)

		case utils.D: r.updateServersCache(s.Service ,false ,hostAndPort)

		}
	})
	log.Println("subscribe over")
}

func (r *goRpc) updateServersCache(serviceName string ,addOrDel bool,host string){
	if addOrDel{
		r.lock <- 1
		for _ , v := range r.serversCache[serviceName]{
			if v == host{
				return //已缓存
			}
		}
		r.serversCache[serviceName] = append(r.serversCache[serviceName],host)
		<- r.lock
	}else {
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


