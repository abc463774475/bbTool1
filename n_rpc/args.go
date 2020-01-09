package n_rpc

import "reflect"

type ArgsInfo struct {
	Name    string     // 函数名字
	C       chan error // channel
	BackErr error      // 返回错误信息

	Param []reflect.Value // 参数列表
}

func createArgsInfo() *ArgsInfo {
	p := &ArgsInfo{}
	p.Param = make([]reflect.Value, 0)
	return p
}

var errInfo = reflect.TypeOf((*error)(nil)).Elem()

type rpc_interface interface {
	GetChan() chan *ArgsInfo
}
