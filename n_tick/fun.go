package n_tick

import (
	"bbTool/n_log"
	"time"
	"runtime/debug"
)

func (self *TickControl) Add(interval int, period time.Duration, f func(i interface{})) int {
	self.iCount++
	st := &BaseInfo{self.iCount, interval, period, false, }

	icount := self.iCount
	time.AfterFunc(period, func() {
		// todo 这儿有bug 哈？？？知道否 ？？？ bug？？  close channel  知道 ？？？
		defer func() {
			if err := recover() ; err != nil {
				n_log.Erro("have fun erro  %v   %v",err)
				debug.PrintStack()
			}
		}()

		if self.C != nil{
			self.C <- icount
		}
	})

	self.data[st.index] = st
	self.fun[st.index] = f

	return self.iCount
}

func (self *TickControl) Del(i int) {
	{
		_, ok := self.data[i]
		if !ok {
			n_log.Erro("cur id is dele", i)
		}else {
			delete(self.data, i)
		}
	}
	{
		_, ok := self.fun[i]
		if !ok {
			n_log.Erro("cur id is dele", i)
		}else {
			delete(self.fun, i)
		}
	}
}

func Handle(tick *TickControl, id int) {
	info := tick.data[id]
	if info == nil {
//		n_log.Debug("tick err  tick have delete  %v",id)
		return
	}

	defer func() {
		delete(tick.data, info.index)
		delete(tick.fun, info.index)
	}()

	if _, ok := tick.data[info.index]; !ok {
		n_log.Info("cur tine is del", info.index)
		return
	}
	df := tick.fun[info.index]
	if df == nil {
		n_log.Info("cur tine is del", info.index)
		return
	}

	df(1)
	// 删除
	if info.totalTimes == -1 {
		tick.Add(info.totalTimes, info.oneTimeLimit, df)
		return
	} else {
		info.totalTimes--
		if info.totalTimes > 0 {
			tick.Add(info.totalTimes, info.oneTimeLimit, df)
		}
	}
}
