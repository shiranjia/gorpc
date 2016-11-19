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
	注册监听事件
	 */
	AddListener(path string  , cancel <- chan int, handler func(*client.Response))

	/**
	根据路径得到key
	 */
	path2key(path string) string

	/**
	根据key得到路径
	 */
	key2path(key string) string
}
