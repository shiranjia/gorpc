package register

import (
	"testing"
	"github.com/coreos/etcd/client"
	"time"
)

func getRegister() Register  {
	r := CreateEtcdRegister("http://192.168.146.147:2379")
	return r
}

func TestEtcdRegister_Get(t *testing.T)  {
	t.Log("Get test")
	r := getRegister()
	t.Log(r.Get("/api.Test"))
}

func TestEtcdRegister_GetChildren(t *testing.T) {
	t.Log("GetChildren test")
	r := getRegister()
	t.Log(r.GetChildren("/"))
}

func TestEtcdRegister_Set(t *testing.T)  {
	t.Log("Set test")
	r := getRegister()
	t.Log(r.Set("/foo/bar","123ee"))
	t.Log(r.Get("/foo/bar"))
}

func TestEtcdRegister_Delete(t *testing.T)  {
	t.Log("Delete test")
	r := getRegister()
	t.Log(r.Delete("/gorpc"))
	//t.Log(r.Delete("/"))
	t.Log(r.GetChildren("/"))
}

func TestEtcdRegister_AddListener(t *testing.T) {
	t.Log("AddListener test")
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

