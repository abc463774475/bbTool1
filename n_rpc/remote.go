package n_rpc

import (
	"bbTool/guid"
	"bbTool/n_log"
	"bbTool/n_routine"
	"net/rpc"
	"sync"
	"time"
)

type poolBaseData struct {
	*rpc.Client            // 用户
	l           sync.Mutex // 锁

	startTime int64 // 开始时间
	endTime   int64 // 结束时间

	addr  string // 地址族
	index string // 唯一编号
}

func CreateBaseData(addr string) *poolBaseData {
	p := &poolBaseData{}
	p.startTime = 0
	p.endTime = 0

	p.addr = addr
	p.index = guid.GetId()

	return p
}

func (self *poolBaseData) Call(fName string, args interface{}, reply interface{}) error {
	// 分配的时候   枷锁
	defer func() {
		G_RemotePool.AddFreeClient(self.addr, self.index, self)
	}()

	s1 := time.Now()
	ctmp := make(chan *rpc.Call, 100)
	n_routine.RoutineFun(self.Go, fName, args, reply, ctmp)

	select {
	case c := <-ctmp:
		end := time.Now().Sub(s1)
		if end >= 1*time.Second {
			n_log.Erro("rpc handle too long ", end, fName)
		}
		return c.Error
	case <-time.After(3 * time.Second):
		n_log.Erro("exe block   %v totalTime %v", fName, time.Now().Sub(s1))
	}

	end := time.Now().Sub(s1)
	if end >= 1*time.Second {
		n_log.Erro("f %v handle rpc very long", fName, end)
	}
	return nil
}
