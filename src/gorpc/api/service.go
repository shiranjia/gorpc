package api

type Facade struct {
	Service string 		`服务全名称`
	Method  string 		`方法`
	Args 	interface{}	`参数`
	Response  interface{}	`返回对象`
	Protocol  string	`协议类型`
}
