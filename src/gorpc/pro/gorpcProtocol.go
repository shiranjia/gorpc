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
tcp transfer gob protocol
 */
func NewRPCServer(service []interface{}){
	for _ ,s := range service {
		log.Println("register rpc service:",reflect.TypeOf(s).String())
		rpc.Register(s)
	}
	listener,err := net.Listen("tcp",":1234")
	utils.CheckErr("gorpcProtocol.NewServer",err)
	go func(l net.Listener){
		for {
			conn,err := l.Accept()
			utils.CheckErr("gorpcProtocol.listen",err)
			rpc.ServeConn(conn)
		}
	}(listener)
}

/**
tcp transfer gob protocol client
 */
func NewRPCClient(host string) *rpc.Client{
	client,err := rpc.Dial("tcp" , host)
	utils.CheckErr("gorpcProtocol.NewClient",err)
	return client
}

/**
http transfer gob protocol
 */
func NewHTTPServer(service []interface{}){
	for _ ,s := range service {
		log.Println("register http service:",reflect.TypeOf(s).String())
		rpc.Register(s)
	}
	rpc.HandleHTTP()
	err := http.ListenAndServe(":1235",nil)
	utils.CheckErr("gorpcProtocol.NewHTTPServer",err)
}

/**
http transfer gob protocol client
 */
func NewHTTPClient(host string) *rpc.Client{
	client,err := rpc.DialHTTP("tcp" , host)
	utils.CheckErr("gorpcProtocol.NewHTTPClient",err)
	return client
}

/**
tcp transfer json protocol
 */
func NewJSONServer(service []interface{})  {
	lis, err := net.Listen("tcp", ":1236")
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

/**
tcp transfer json protocol client
 */
func NewJSONClient(host string) *rpc.Client{
	client, err := jsonrpc.Dial("tcp", host)
	utils.CheckErr("gorpcProtocol.NewJSONClient",err)
	return client
}

/**
tcp transfer json2 protocol
 */
func NewJSON2Server(service []interface{}){
	for _,s := range service{
		log.Println("register json2rpc service:",reflect.TypeOf(s).String())
		err := rpc.Register(s)
		utils.CheckErr("gorpcProtocol.NewJSON2Server.Register",err)
	}
	listener, err := net.Listen("tcp", ":1237")
	utils.CheckErr("gorpcProtocol.NewJSON2Server",err)
	go func(lis net.Listener){
		for  {
			con,err := lis.Accept()
			utils.CheckErr("gorpcProtocol.NewJSON2Server.Accept",err)
			go jsonrpc2.ServeConn(con)
		}
	}(listener)
}

/**
tcp transfer json2 protocol client
 */
func NewJSON2Client(host string) *jsonrpc2.Client  {
	client,err := jsonrpc2.Dial("tcp",host)
	utils.CheckErr("gorpcProtocol.NewJSON2Client",err)
	return client
}

/**
http transfer json protocol
 */
func NewHttpJson2rpcServer(service [] interface{}) {
	server := rpc.NewServer()
	for _,s := range service{
		log.Println("register sttpJson2rpc service:",reflect.TypeOf(s).String())
		err := server.Register(s)
		utils.CheckErr("gorpcProtocol.NewHttpJson2rpcServer.Register",err)
	}
	// Server provide a HTTP transport on /rpc endpoint.
	http.Handle("/rpc", jsonrpc2.HTTPHandler(server))
	lnHTTP, err := net.Listen("tcp", ":1238")
	utils.CheckErr("gorpcProtocol.NewHttpJson2rpcServer",err)
	go http.Serve(lnHTTP, nil)
}

/**
http transfer json2 protocol
 */
func NewHttpJson2rpcClient(host string) *jsonrpc2.Client {
	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://" + host + "/rpc")
	return clientHTTP
}