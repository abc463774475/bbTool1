package n_rpc

import (
	"bbTool/guid"
	"bbTool/n_log"
	"bbTool/n_routine"
	"errors"
	"fmt"
	"reflect"
)

type LocalRpcData struct {
	C chan *ArgsInfo
}

func (self *LocalRpcData) GetChan() chan *ArgsInfo {
	return self.C
}

func CreateLocalRpcData() *LocalRpcData {
	p := &LocalRpcData{}
	p.Init()
	return p
}

func (self *LocalRpcData) Init() {
	self.C = make(chan *ArgsInfo, 100)
}

func (self LocalRpcData) Del() {
	close(self.C)
	self.C = nil
}

func (self *LocalRpcData)Call_notBack(name string, v ...interface{})  {
	v = append([]interface{}{name},v...)
	n_routine.RoutineFun(self.Call,v...)
	//go self.Call(name,v...)
}

func (self *LocalRpcData) Call(name string, v ...interface{}) error {

	key := guid.GetId()
	n_log.OnlyFile("  rpc start  %v name %v", key,name)

	ctmp := make(chan error, 1)

	defer func() {
		n_log.OnlyFile("rpc end  %v", key)
		close(ctmp)
	}()

	n_routine.RoutineFun(func() {
		c1 := createArgsInfo()
		c1.Name = name

		for _, v1 := range v {
			c1.Param = append(c1.Param, reflect.ValueOf(v1))
		}

		c1.C = ctmp
		self.C <- c1
	})

	err, _ := <-ctmp

	return err
}

func Loop(c rpc_interface) {
	n_log.Info("start  rpc  ")
	defer n_log.Info("cur end")

	for {
		select {
		case data, ok := <-c.GetChan():
			n_log.Info("data  %v", data)
			if !ok {
				n_log.Info("tablemgrh exit")
				return
			}

			HandleRpc(c, data)
		}
	}
}

func HandleRpc(info rpc_interface, data *ArgsInfo) {
	r := reflect.ValueOf(info)

	mFun := r.MethodByName(data.Name)

	n_log.Info("data %v  %v  %v",data.Name,len(data.Param),data.Param)
	if mFun.Kind() == 0 {
		n_log.Erro("not have this f %v  ", data.Name)
	} else {
		n_log.Info("call %v", data.Param[0])
		b := mFun.Call(data.Param)
		if len(b) != 1 {
			data.BackErr = errors.New(fmt.Sprintf("back len erro %d", len(b)))
		} else {
			// 回传error
			if b[0].Type() != errInfo {
				data.BackErr = errors.New(fmt.Sprintf("retruen type not right  %v", b[0].Type()))
			} else if b[0].IsNil() {

			} else {
				data.BackErr = b[0].Interface().(error)
			}
		}
	}

	data.C <- data.BackErr
}
