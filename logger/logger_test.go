package logger

import (
	"os"
	"testing"
	"time"
)

func Test_GetWriter(t *testing.T) {

	out := getWriter("all", os.Stdout, os.Stdout)
	out.Write([]byte("hello world"))
	time.Sleep(100000)
	t.Log("good")
}

func Test_Init(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Error(err)
		}

	}()
	Init("debug", "all")
	Debugln("debug_______test")
	Errorln("error_______test")
}
