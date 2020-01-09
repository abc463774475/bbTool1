package n_rpc

import (
	"bbTool/n_log"
	"testing"
	"net/rpc"
	"net"
	"bbTool/n_routine"
	"time"
)

func TestCallRemote(t *testing.T) {

	n_log.Erro("222222222")
}

type Proxy_1 int

func (self *Proxy_1)F1(args []byte,reply *[]byte)  error {
	n_log.Debug("2222222")
	*reply = []byte("gggggggg")
	return nil
}

type Proxy_2 int

func (self *Proxy_2)F1(args []byte,reply *[]byte)  error {
	n_log.Debug("555555")
	*reply = []byte("hhhxxx")
	return nil
}

func Test11111111(t *testing.T)  {
	newServer := rpc.NewServer()

	newServer.Register(new(Proxy_1))
	newServer.Register(new(Proxy_2))

	addr := "127.0.0.1:1234"
	l, e := net.Listen("tcp", addr) // any available address
	if e != nil {
		n_log.Panic("err  %v",e)
	}
	//go newServer.Accept(l)

	n_log.Info("初始化  %v rpc 成功", addr)

	n_routine.RoutineFun(newServer.Accept, l)

	for {
		time.Sleep(10*time.Second)
	}
}

func TestClientok(t *testing.T)  {
	addr := "127.0.0.1:1234"

	var reply []byte
	err := CallRemote(addr, "Proxy_2.F1",[]byte("xxxxx"),&reply)

	n_log.Debug("errr  %v",err)
	n_log.Debug("reply  %v",string(reply))
}