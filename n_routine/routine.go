package n_routine

import (
	"reflect"
	"bbTool/n_log"
	"runtime/debug"
)

//  函数后面跟参数列表   ok？
func RoutineFun(f interface{}, i ...interface{}) {
	go func(f interface{}, i ...interface{}) {
		defer func() {
			if err := recover() ; err != nil {
				r := reflect.TypeOf(f)
				n_log.Erro("have fun erro  %v   %v",err,r)
				debug.PrintStack()
			}
		}()

		j := reflect.ValueOf(f)

		param := []reflect.Value{}
		for _, v := range i {
			param = append(param, reflect.ValueOf(v))
		}
		j.Call(param)
	}(f, i...)
}
