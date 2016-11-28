#  gorpc
## go语言分布式服务总线
---------------------------------------------------------------------------------------------------------------------------------------------
##注册中心使用etcd
主要功能有服务自动发现，负载均衡，故障转移，监控
序列化协议支持gob,json,json-rpc
通信支持tcp，http

###例子
注册服务：
    etcdUrl := "http://192.168.146.147:2379"
	rpc := NewGoRpc(etcdUrl)
	rpc.RegisterServer(
		service.Service{&Test{},utils.PROTOCOL_RPC},
	)
消费服务：
    etcdUrl := "http://192.168.146.147:2379"
	goRpc := NewGoRpc(etcdUrl)
	f := Facade{
		Service:"api.Test",
		Method:"Tostring",
		Args:Request{"ttt protocol rpc"},
		Response:&Response{},
		Protocol:utils.PROTOCOL_RPC,
	}
	goRpc.Call(f)
	t.Log(f.Response)

----------------------------------------------------------------------------------------------------------------------------------------------
部署etcd:
$ etcd --name infra0 --initial-advertise-peer-urls http://192.168.146.147:2380 \
  --listen-peer-urls http://192.168.146.147:2380 \
  --listen-client-urls http://192.168.146.147:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://192.168.146.147:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster infra0=http://192.168.146.147:2380 \
  --initial-cluster-state new
