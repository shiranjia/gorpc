package api

import (
	"gorpc/pro"
	"gorpc/server"
	"gorpc/register"
	_ "math/rand"
	"gorpc/utils"
	"reflect"
	"log"
	"strings"
	"net/rpc"
	"math/rand"
)

type goRpc struct {
	serversCache map[string][]*server.Provider `服务列表`
	point server.Provider
	register.Register
}

func NewGoRpc(host string) *goRpc{
	g := &goRpc{}
	g.serversCache = make(map[string][]*server.Provider)
	g.Register = register.CreateEtcdRegister(host)
	return  g
}

func (r *goRpc) RegisterServer(service interface{})  {
	t :=reflect.TypeOf(service)
	serviceName := t.String()
	serviceName = strings.Replace(serviceName,"*","",-1)
	log.Println(serviceName)
	//s := buildService(serviceName + "/" + utils.Ip())
	r.Register.Set(serviceName + "/" + utils.Ip() +":1234" , "")
	pro.NewServer(service)
}

func (r *goRpc) Call(s Facade) error  {
	nodes,err := r.Register.GetChildren(s.Service)
	utils.CheckErr(err)
	log.Println(nodes)
	client := pro.NewClient(nodes[0].Key)
	ss := strings.Split(s.Service,".")
	m := ss[len(ss) - 1] + "." + s.Method
	log.Println("call method :",m)
	return client.Call(m,s.Args,s.Response)
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
	pro.NewHTTPServer(services)
}

func (r *goRpc) CallHTTP(s Facade) error  {
	hosts := r.serversCache[s.Service]
	var client *rpc.Client
	var method string
	if hosts == nil || len(hosts)==0{
		nodes,err := r.Register.GetChildren(s.Service)
		utils.CheckErr(err)
		log.Println(nodes)
		client = pro.NewHTTPClient(nodes[0].Key)
		ss := strings.Split(s.Service,".")
		method = ss[len(ss) - 1] + "." + s.Method
		log.Println("call method :",method)
		p := &server.Provider{nodes[0].Key,method}
		r.serversCache[s.Service] = append(r.serversCache[s.Service],p)
	}else {
		prov := hosts[rand.Int() % len(hosts)]
		client = pro.NewHTTPClient(prov.Host)
		method = prov.Method
		log.Println("point cache host:", prov.Host,"method:",prov.Method)
	}

	return client.Call(method,s.Args,s.Response)
}

func buildService(service string) server.Provider{
	s := server.Provider{}
	s.Host = utils.Ip() + "1234"
	return s
}



type Request struct {
	Body interface{}
}

type Response struct {
	Body interface{}
}


