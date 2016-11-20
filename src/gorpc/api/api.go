package api

import (
	"gorpc/pro"
	"gorpc/server"
	"gorpc/register"
	"math/rand"
	"strconv"
	"gorpc/utils"
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

func (r *goRpc) RegisterServer(serviceName string,service interface{})  {
	s := buildService(serviceName)
	r.Register.Set(s.ServiceName,s.Host + strconv.Itoa(s.Port))
	pro.NewServer(service)
}

func (r *goRpc) ProxyServer()  {

}

func (r *goRpc) Call(s service) error  {
	servers := r.servers[s.Service + "." + s.Method]
	server := servers[rand.Int()%len(servers)]
	client := pro.NewClient(server.Host)
	return client.Call(s.Method,s.Args,s.Response)
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


