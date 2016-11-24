package utils

import (
	"testing"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

func TestRoundSelector(t *testing.T) {
	i := RoundSelector(2)
	t.Log(i)
}

func TestTry(t *testing.T) {
	Try(func(){
		panic(errors.New("test"))
	},func(err interface{}){
		t.Log("catch code block")
	})
}