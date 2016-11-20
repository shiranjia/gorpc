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
)

type goRpc struct {
	servers map[string][]server.Provider `服务列表`
	point server.Provider
	register.Register
}

func NewGoRpc(host string) *goRpc{
	g := &goRpc{}
	g.servers = make(map[string][]server.Provider)
	g.Register = register.CreateEtcdRegister(host)
	return  g
}

func (r *goRpc) RegisterServer(service interface{})  {
	t :=reflect.TypeOf(service)
	serviceName := t.String()
	serviceName = strings.Replace(serviceName,"*","",-1)
	log.Println(serviceName)
	//s := buildService(serviceName + "/" + utils.Ip())
	r.Register.Set(serviceName + "/" + "127.0.0.1" +":7777" , "")
	pro.NewServer(service)
}

func (r *goRpc) ProxyServer()  {

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

func buildService(service string) server.Provider{
	s := server.Provider{}
	s.Port = 8888
	s.Host = utils.Ip()
	return s
}



type Request struct {
	Body interface{}
}

type Response struct {
	Body interface{}
}


