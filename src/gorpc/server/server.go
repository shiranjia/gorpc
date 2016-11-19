package server

/**
服务节点机器信息
 */
type Server struct {
	Host string 	`服务ip`
	Port int	`服务端口`
	Next Server	`下一节点`
}




