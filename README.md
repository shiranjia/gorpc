# gorpc
go语言分布式服务总线

部署etcd:
$ etcd --name infra0 --initial-advertise-peer-urls http://192.168.146.147:2380 \
  --listen-peer-urls http://192.168.146.147:2380 \
  --listen-client-urls http://192.168.146.147:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://192.168.146.147:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster infra0=http://192.168.146.147:2380 \
  --initial-cluster-state new
