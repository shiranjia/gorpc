package api

import (
	"testing"
)

func TestMonitor_GetDate(t *testing.T) {
	etcdUrl := "http://127.0.0.1:2379"
	rpc := NewGoRpc(etcdUrl)
	monitor := rpc.Monitor
	monitor.GetDate()
	t.Log(rpc.Monitor.Service)
}

