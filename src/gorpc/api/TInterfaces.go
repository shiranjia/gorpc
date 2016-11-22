package api

import (
	"log"
	"fmt"
)
//测试使用注册服务
type Test struct {}
func (t *Test) Tostring(req Request,resp *Response)  error {
	log.Println(req.Body)
	resp.Body = fmt.Sprint(req.Body) +",test"
	return nil
}
type Test1 struct {}
func (t *Test1) Tostring(req Request,resp *Response)  error {
	log.Println(req.Body)
	resp.Body = fmt.Sprint(req.Body) +",test1"
	return nil
}
