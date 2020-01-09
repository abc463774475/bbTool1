package n_rpc

import (
	"net/rpc"
	"net"
	"bbTool/n_log"
	"bbTool/n_routine"
)

func Init_Rpc(addr string, i... interface{}) {
	newServer := rpc.NewServer()
	for _,v := range i {
		newServer.Register(v)
	}

	l, e := net.Listen("tcp", addr) // any available address
	if e != nil {
		n_log.Panic("",e)
	}
	//go newServer.Accept(l)

	n_log.Info("初始化  %v rpc 成功", addr)
	n_routine.RoutineFun(newServer.Accept, l)
}
