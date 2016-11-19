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

	rootPath string `根目录`

	separator string `目录分隔符`

	client client.KeysAPI
}

func (r * etcdRegister) Connect()  {
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
		r.rootPath = "/gorpc"
		r.separator = "/"
	},func(v interface{}){
		log.Fatal("create etcd client err:",v)
	})
}

func (r *etcdRegister) Set(path string,value string) error  {
	res ,err := r.client.Set(context.Background(),path,value,nil)
	err = c.CheckErr(err)
	_ = res
	return err
}

func (r *etcdRegister) Get(path string) (Node,error) {
	res,err := r.client.Get(context.Background(),path,nil)
	err = c.CheckErr(err)
	var v Node
	v.Key = r.path2key(path)
	v.Path = path
	if res!= nil && res.Node!=nil{
		v.Value = res.Node.Value
	}
	return v,err
}

func (r *etcdRegister) GetChildren(path string) ([]Node,error){
	res,err := r.client.Get(context.Background(),path,&client.GetOptions{true,false,true})
	err = c.CheckErr(err)
	var nodes []Node
	if res!= nil && res.Node!=nil{
		childs := res.Node.Nodes
		for _ ,c := range childs{
			var n Node
			n.Key = r.path2key(c.Key)
			n.Path = c.Key
			n.Value = c.Value
			nodes = append(nodes , n)
		}
	}
	return nodes,err
}

func (r *etcdRegister) Delete(path string) error  {
	res,err := r.client.Delete(context.Background(),path,nil)
	err = c.CheckErr(err)
	_ = res
	return err
}

func (r *etcdRegister) AddListener(path string , cancel <- chan int,
					handler func(*client.Response))  {
	watcher := r.client.Watcher(path,&client.WatcherOptions{0,true})
	go func(client.Watcher){
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
	}(watcher)
}

func (r *etcdRegister) path2key(path string) string{
	ps := strings.Split(path,r.separator)
	l := len(ps)
	key := ps[0]
	if l > 0 {
		key = ps[l - 1]
	}
	return key
}

func(r *etcdRegister) key2path(key string) string{
	return key
}
