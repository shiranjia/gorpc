package pro

import (
	"net"
	"gorpc/utils"
	"net/rpc"
	"log"
	"net/http"
	"reflect"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"net/rpc/jsonrpc"
)

/**
gob
 */
func NewRPCServer(service []interface{}){
	for _ ,s := range service {
		log.Println("register rpc service:",reflect.TypeOf(s).String())
		rpc.Register(s)
	}
	listener,err := net.Listen("tcp",":1234")
	utils.CheckErr("gorpcProtocol.NewServer",err)
	go listen(listener)
}

func listen(l net.Listener){
	for {
		conn,err := l.Accept()
		utils.CheckErr("gorpcProtocol.listen",err)
		rpc.ServeConn(conn)
	}
}

func NewRPCClient(host string) *rpc.Client{
	client,err := rpc.Dial("tcp" , host)
	utils.CheckErr("gorpcProtocol.NewClient",err)
	return client
}

func NewHTTPServer(service []interface{}){
	for _ ,s := range service {
		log.Println("register http service:",reflect.TypeOf(s).String())
		rpc.Register(s)
	}
	rpc.HandleHTTP()
	err := http.ListenAndServe(":1234",nil)
	utils.CheckErr("gorpcProtocol.NewHTTPServer",err)
}

func NewHTTPClient(host string) *rpc.Client{
	client,err := rpc.DialHTTP("tcp" , host)
	utils.CheckErr("gorpcProtocol.NewHTTPClient",err)
	return client
}

func NewJSONServer(service []interface{})  {
	lis, err := net.Listen("tcp", ":1234")
	utils.CheckErr("gorpcProtocol.NewJSONServer",err)
	srv := rpc.NewServer()
	for _,s := range service {
		log.Println("register json service:",reflect.TypeOf(s).String())
		err := srv.Register(s)
		utils.CheckErr("gorpcProtocol.NewJSONServer.Register",err)
	}
	go func(l net.Listener,ser *rpc.Server){
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Fatalf("lis.Accept(): %v\n", err)
			}
			go srv.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}(lis,srv)
}

func NewJSONClient(host string) *rpc.Client{
	client, err := jsonrpc.Dial("tcp", host)
	utils.CheckErr("gorpcProtocol.NewJSONClient",err)
	return client
}

func NewJSON2Server(service []interface{}){
	for _,s := range service{
		log.Println("register json2rpc service:",reflect.TypeOf(s).String())
		err := rpc.Register(s)
		utils.CheckErr("gorpcProtocol.NewJSON2Server.Register",err)
	}
	listener, err := net.Listen("tcp", ":1234")
	utils.CheckErr("gorpcProtocol.NewJSON2Server",err)
	go func(lis net.Listener){
		for  {
			con,err := lis.Accept()
			utils.CheckErr("gorpcProtocol.NewJSON2Server.Accept",err)
			go jsonrpc2.ServeConn(con)
		}
	}(listener)
}

func NewJSON2Client(host string) *jsonrpc2.Client  {
	client,err := jsonrpc2.Dial("tcp",host)
	utils.CheckErr("gorpcProtocol.NewJSON2Client",err)
	return client
}

func NewHttpJson2rpcServer(service [] interface{}) {
	server := rpc.NewServer()
	for _,s := range service{
		log.Println("register sttpJson2rpc service:",reflect.TypeOf(s).String())
		err := server.Register(s)
		utils.CheckErr("gorpcProtocol.NewHttpJson2rpcServer.Register",err)
	}
	// Server provide a HTTP transport on /rpc endpoint.
	http.Handle("/rpc", jsonrpc2.HTTPHandler(server))
	lnHTTP, err := net.Listen("tcp", ":1235")
	utils.CheckErr("gorpcProtocol.NewHttpJson2rpcServer",err)
	go http.Serve(lnHTTP, nil)
}

func NewHttpJson2rpcClient(host string) *jsonrpc2.Client {
	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://" + host + "/rpc")
	return clientHTTP
}