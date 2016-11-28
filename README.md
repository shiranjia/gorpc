##  gorpc
###go语言分布式服务总线
----------------------------------------------------------------------------------------------------------------------------------------
##### * 注册中心使用etcd</br>
##### * 主要功能有服务自动发现，负载均衡，故障转移，监控</br>
##### * 序列化协议支持gob，json，json-rpc</br>
##### * 通信支持tcp，http</br>

###例子
#####定义服务：</br>
#####type Test struct {} </br>
#####func (t *Test) Tostring(req Request,resp *Response)  error {</br>
#####	log.Println(req.Body)</br>
#####	resp.Body = fmt.Sprint(req.Body) +",test"</br>
#####	return nil</br>
#####}</br>
#####注册服务：</br>
#####etcdUrl := "http://192.168.146.147:2379"</br>
#####rpc := NewGoRpc(etcdUrl)</br>
#####rpc.RegisterServer(</br>
#####service.Service{&Test{},utils.PROTOCOL_RPC},</br>
#####)</br>

#####消费服务：</br>
#####   etcdUrl := "http://192.168.146.147:2379"</br>
#####	goRpc := NewGoRpc(etcdUrl)</br>
#####	f := Facade{</br>
#####		Service:"api.Test",</br>
#####		Method:"Tostring",</br>
#####		Args:Request{"ttt protocol rpc"},</br>
#####		Response:&Response{},</br>
#####		Protocol:utils.PROTOCOL_RPC,</br>
#####	}</br>
#####	goRpc.Call(f)</br>
#####	t.Log(f.Response)</br>

------------------------------------------------------------------------------------------------------------------------------------------
##### static模式部署etcd:
#####$ etcd --name infra0 --initial-advertise-peer-urls http://192.168.146.147:2380 \
#####  --listen-peer-urls http://192.168.146.147:2380 \
#####  --listen-client-urls http://192.168.146.147:2379,http://127.0.0.1:2379 \
#####  --advertise-client-urls http://192.168.146.147:2379 \
#####  --initial-cluster-token etcd-cluster-1 \
#####  --initial-cluster infra0=http://192.168.146.147:2380 \
#####  --initial-cluster-state new
