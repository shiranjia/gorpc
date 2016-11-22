package register

import (
	"github.com/coreos/etcd/client"
	"context"
	c "gorpc/utils"
	"log"
	"time"
	"strings"
)

type etcdRegister struct {

	/**
	例：http://127.0.0.1:2379;http://127.0.0.2:2379
	 */
	host string `服务地址`

	client client.KeysAPI

	hosts []string	`保存服务列表，定时推送etcd`

	updateInterval time.Duration	`心跳时间`
}

//构造函数
func CreateEtcdRegister(host string) Register  {
	var r Register
	c.Try(func(){
		etcd := &etcdRegister{
			host : host,
			updateInterval : 5 * time.Second,
			hosts : []string{},
		}
		r = etcd
		r.connect()
	}, func(err interface{}) {
		log.Fatalln("create register err:",err)
	})
	return r
}

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
	log.Println("set path:",c.Key2path(path))
	res ,err := r.client.Set(context.Background(),c.Key2path(path),value,
		&client.SetOptions{
			TTL : r.updateInterval + 10 * time.Second,
			PrevExist : client.PrevIgnore,
		})
	err = c.CheckErr("etcdRegister.Set",err)
	_ = res
	r.hosts = append(r.hosts,c.Key2path(path))
	return err
}

func (r *etcdRegister) TimeTicker()  {
	var ticker *time.Ticker = time.NewTicker(r.updateInterval)
	go func(){
		for range ticker.C{
			for _ , h := range r.hosts{
				res ,err := r.client.Set(context.Background(),h,"1",
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

func (r *etcdRegister) Subscribe(path string , cancel <- chan int,
					handler func(cli *client.Response))  {
	watcher := r.client.Watcher(path,&client.WatcherOptions{0,true})
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



