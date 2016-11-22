package register

import (
	"github.com/coreos/etcd/client"
)

type Register interface {

	/**
	建立链接
	 */
	connect()

	/**
	设置元素
	 */
	Set(path string,value string) error

	/**
	获得元素
	 */
	Get(path string) (Node,error)

	/**
	获得元素
	 */
	GetChildren(path string) ([]Node,error)

	/**
	删除元素
	 */
	Delete(path string) error

	/**
	订阅服务变化
	 */
	Subscribe(path string  , cancel <- chan int, handler func(*client.Response))

	/**
	维持心跳  etcd需要自己维持心跳保持临时节点有效
	 */
	TimeTicker()

}
