package pro

import (
	"testing"
	"fmt"
	"strconv"
)

func TestNewRPCServer(t *testing.T) {
	t.Log("test tNewRPCServer")
	NewRPCServer([]interface{}{&ExampleSvc{}})
	w := make(chan int)
	<- w
}

func TestNewRPCClient(t *testing.T) {
	t.Log("test tNewHTTPClient")
	host := "127.0.0.1:1234"
	client := NewRPCClient(host)
	defer client.Close()
	reply := &Response{}
	err := client.Call("ExampleSvc.Sum", [2]int{3, 5}, reply)
	if err!=nil{
		t.Log(err)
	}
	fmt.Printf("Sum(3,5)=%s\n", reply)
}

func TestNewHTTPServer(t *testing.T) {
	t.Log("test tNewHTTPServer")
	NewHTTPServer([]interface{}{&ExampleSvc{}})
	w := make(chan int)
	<- w
}

func TestNewHTTPClient(t *testing.T) {
	t.Log("test tNewHTTPClient")
	host := "127.0.0.1:1234"
	client := NewHTTPClient(host)
	defer client.Close()
	reply := &Response{}
	err := client.Call("ExampleSvc.Sum", [2]int{3, 5}, reply)
	if err!=nil{
		t.Log(err)
	}
	fmt.Printf("Sum(3,5)=%s\n", reply)
}

func TestNewJSONServer(t *testing.T) {
	t.Log("test tNewJSONServer")
	NewJSONServer([]interface{}{&ExampleSvc{}})
	w := make(chan int)
	<- w
}

func TestNewJSONClient(t *testing.T) {
	t.Log("test tNewJSONClient")
	host := "127.0.0.1:1234"
	client := NewJSONClient(host)
	defer client.Close()
	reply := &Response{}
	err := client.Call("ExampleSvc.Sum", [2]int{3, 5}, reply)
	if err!=nil{
		t.Log(err)
	}
	fmt.Printf("Sum(3,5)=%s\n", reply)
}

func TestNewJSON2Server(t *testing.T) {
	t.Log("test tNewJSON2Server")
	NewJSON2Server([]interface{}{&ExampleSvc{}})
	w := make(chan int)
	<- w
}

func TestNewJSON2Client(t *testing.T) {
	t.Log("test tNewJSON2Client")
	// Server provide a TCP transport.
	host := "127.0.0.1:1234"
	client := NewJSON2Client(host)
	defer client.Close()
	reply := &Response{}
	err := client.Call("ExampleSvc.Sum", [2]int{3, 5}, reply)
	if err!=nil{
		t.Log(err)
	}
	fmt.Printf("Sum(3,5)=%s\n", reply)
}

type ExampleSvc struct{}
func (*ExampleSvc) Sum(vals [2]int, res *Response) error {
	fmt.Printf("arg1=%d,arg2=%d",vals[0],vals[1])
	res.Body = strconv.Itoa(vals[0]) + strconv.Itoa(vals[1])
	return nil
}
type Response struct {
	Body string
}
