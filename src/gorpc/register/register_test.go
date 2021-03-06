package register

import (
	"testing"
	"github.com/coreos/etcd/client"
	"time"
	"strings"
)

func getRegister() Register  {
	r := CreateEtcdRegister("http://127.0.0.1:2379")
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
	t.Log(r.GetChildren("api.Test"))
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
	//t.Log(r.Delete("/gorpc"))
	t.Log(r.Delete("/"))
	t.Log(r.GetChildren("/"))
}

func TestEtcdRegister_AddListener(t *testing.T) {
	t.Log("AddListener test")
	r := getRegister()
	r.Subscribe("/gorpc/api.Test1",make(chan int), func(c *client.Response) {
		path := strings.Split(c.Node.Key,"/")
		t.Log(c.Action," 》",path[len(path)-1])
	})
	r.Set("api.Test1/127.0.0.1:1234","asd")
	r.Delete("api.Test1/127.0.0.1:1234")
	time.Sleep(3 * time.Second)
}

