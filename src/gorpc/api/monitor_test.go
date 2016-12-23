package api

import (
	"testing"
)

func TestMonitor_GetDate(t *testing.T) {
	etcdUrl := "http://127.0.0.1:2379"
	rpc := NewGoRpc(etcdUrl)
	monitor := rpc.Monitor
	monitor.GetDate()
	for k,v := range rpc.Monitor.Service {
		t.Log(k,"->","Providers:",v.Provider,"Consumers:",v.Consumer)
	}

}

