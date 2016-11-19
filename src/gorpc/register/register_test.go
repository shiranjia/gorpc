package register

import (
	"testing"
	"github.com/coreos/etcd/client"
	"time"
)

func getRegister() Register  {
	var r Register
	etcd := &etcdRegister{}
	etcd.host = "http://127.0.0.1:2379"
	r = etcd
	r.Connect()
	return r
}

func TestEtcdRegister_Get(t *testing.T)  {
	t.Log("get test")
	r := getRegister()
	t.Log(r.Get("/foo"))
}

func TestEtcdRegister_GetChildren(t *testing.T) {
	t.Log("GetChildren test")
	r := getRegister()
	t.Log(r.GetChildren("/foo"))
}

func TestEtcdRegister_Set(t *testing.T)  {
	t.Log("set test")
	r := getRegister()
	t.Log(r.Set("/foo/bar","123eee"))
	t.Log(r.Get("/foo/bar"))
}

func TestEtcdRegister_Delete(t *testing.T)  {
	t.Log("set delete")
	r := getRegister()
	t.Log(r.Delete("/foo/bar"))
	t.Log(r.Delete("/foo"))
	t.Log(r.Get("/foo"))
}

func TestEtcdRegister_AddListener(t *testing.T) {
	t.Log("set AddListener")
	r := getRegister()
	r.AddListener("/foo",make(chan int), func(c *client.Response) {
		t.Log(c.Action," 》",c.Node)
	})
	r.Set("/foo","asd")
	r.Set("/foo","123")
	r.Set("/foo/bar","123")
	r.Delete("/foo")
	time.Sleep(3 * time.Second)
}
