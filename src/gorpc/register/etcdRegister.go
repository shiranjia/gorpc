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
package register

import (
	"github.com/coreos/etcd/client"
	"context"
	c "gorpc/utils"
	"log"
	"time"
	"strings"
	"gorpc/service"
	"strconv"
)

type etcdRegister struct {

	/**
	例：http://127.0.0.1:2379;http://127.0.0.2:2379
	 */
	host string `服务地址`

	client client.KeysAPI

	hosts []*service.Provider	`保存服务列表，定时推送etcd`

	updateInterval time.Duration	`心跳时间`
}

//构造函数
func CreateEtcdRegister(host string) Register  {
	var r Register
	c.Try(func(){
		etcd := &etcdRegister{
			host : host,
			updateInterval : 5 * time.Second,
			hosts : []*service.Provider{},
		}
		r = etcd
		r.connect()
	}, func(err interface{}) {
		log.Fatalln("create register err:",err)
	})
	return r
}

/**
初始化链接
 */
func (r * etcdRegister) connect()  {
	c.Try(func(){
		cfg := client.Config{
			Endpoints:               strings.Split(r.host,";"),
			Transport:               client.DefaultTransport,
			// set timeout per request to fail fast when the target endpoint is unavailable
			HeaderTimeoutPerRequest: time.Second,
		}
		c, err := client.New(cfg)
		if err != nil {
			log.Fatal(err)
		}
		r.client = client.NewKeysAPI(c)
	},func(v interface{}){
		log.Fatal("create etcd client err:",v)
	})
}

func (r *etcdRegister) Set(path string,value string) error  {
	address := c.Key2path(path)
	log.Println("set path:",address)
	res ,err := r.client.Set(context.Background(),c.Key2path(path),value,
		&client.SetOptions{
			TTL : r.updateInterval + 10 * time.Second,
			PrevExist : client.PrevIgnore,
		})
	err = c.CheckErr("etcdRegister.Set",err)
	_ = res
	r.hosts = append(r.hosts,&service.Provider{
		Address:address,
		Invoke:0,
	})
	return err
}

/**
保持心跳，维持临时节点
 */
func (r *etcdRegister) TimeTicker()  {
	var ticker *time.Ticker = time.NewTicker(r.updateInterval)
	go func(){
		for range ticker.C{
			for _ , h := range r.hosts{
				res ,err := r.client.Set(context.Background(),h.Address,strconv.Itoa(h.Invoke),
					&client.SetOptions{
						TTL : r.updateInterval + 10 * time.Second,
						PrevExist : client.PrevIgnore,
					})
				err = c.CheckErr("etcdRegister.TimeTicker",err)
				_ = res
			}
		}
	}()
}

func (r *etcdRegister) Get(path string) (Node,error) {
	res,err := r.client.Get(context.Background(),path,nil)
	err = c.CheckErr("etcdRegister.Get",err)
	var v Node
	v.Key = c.Path2key(path)
	v.Path = path
	if res!= nil && res.Node!=nil{
		v.Value = res.Node.Value
	}
	return v,err
}

func (r *etcdRegister) GetChildren(path string) ([]Node,error){
	res,err := r.client.Get(context.Background(),c.Key2path(path),&client.GetOptions{true,false,true})
	err = c.CheckErr("etcdRegister.GetChildren",err)
	var nodes []Node
	if res!= nil && res.Node!=nil{
		childs := res.Node.Nodes
		for _ ,child := range childs{
			var n Node
			n.Key = c.Path2key(child.Key)
			n.Path = child.Key
			n.Value = child.Value
			nodes = append(nodes , n)
		}
	}
	return nodes,err
}

func (r *etcdRegister) Delete(path string) error  {
	res,err := r.client.Delete(context.Background(),c.Key2path(path),&client.DeleteOptions{
		Recursive : true,
	})
	err = c.CheckErr("etcdRegister.Delete",err)
	_ = res
	return err
}

/**
订阅路径变更
 */
func (r *etcdRegister) Subscribe(path string , cancel <- chan int,
					handler func(cli *client.Response))  {
	watcher := r.client.Watcher(path,&client.WatcherOptions{0,true})
	log.Println("订阅事件，目录：",path)
	for {
		select {
		case <- cancel:
			log.Println("exit watcher")
			return
		default:{
			res , err := watcher.Next(context.Background())
			if err != nil{
				log.Fatal(err)
			}
			handler(res)
		}
		}
	}
}



