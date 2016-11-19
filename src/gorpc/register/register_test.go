package register

import (
	"testing"
	"github.com/coreos/etcd/client"
	"time"
)

func TestEtcdRegister_Get(t *testing.T)  {
	t.Log("get test")
	var r Register
	etcd := &etcdRegister{}
	etcd.host = "http://127.0.0.1:2379"
	r = etcd
	r.Connect()
	t.Log(r.Get("/foo"))
}

func TestEtcdRegister_GetChildren(t *testing.T) {
	t.Log("GetChildren test")
	var r Register
	etcd := &etcdRegister{}
	etcd.host = "http://127.0.0.1:2379"
	r = etcd
	r.Connect()
	t.Log(r.GetChildren("/foo"))
}

func TestEtcdRegister_Set(t *testing.T)  {
	t.Log("set test")
	var r Register
	etcd := &etcdRegister{}
	etcd.host = "http://127.0.0.1:2379"
	r = etcd
	r.Connect()
	t.Log(r.Set("/foo/bar","123eee"))
	t.Log(r.Get("/foo/bar"))
}

func TestEtcdRegister_Delete(t *testing.T)  {
	t.Log("set delete")
	var r Register
	etcd := &etcdRegister{}
	etcd.host = "http://127.0.0.1:2379"
	r = etcd
	r.Connect()
	t.Log(r.Delete("/foo"))
	t.Log(r.Get("/foo"))
}

func TestEtcdRegister_AddListener(t *testing.T) {
	t.Log("set AddListener")
	var r Register
	etcd := &etcdRegister{}
	etcd.host = "http://127.0.0.1:2379"
	r = etcd
	r.Connect()
	r.AddListener("/foo",make(chan int), func(c *client.Response) {
		t.Log(c.Action," 》",c.Node)
	})
	r.Set("/foo","asd")
	r.Set("/foo","123")
	r.Set("/foo/bar","123")
	r.Delete("/foo")
	time.Sleep(3 * time.Second)
}
