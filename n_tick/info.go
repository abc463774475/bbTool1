package n_tick

import "time"

type TickFun func(i interface{})

type BaseInfo struct {
	index        int           // 标号
	totalTimes   int           // 总共次数
	oneTimeLimit time.Duration // 持续时间
	isClose bool    		   // 是否删除
}

type TickControl struct {
	C      chan int
	iCount int // 计数器  上面的标志

	data map[int]*BaseInfo // 数据
	fun  map[int]TickFun
}

func (self *TickControl) Init() {
	self.C = make(chan int, 100)
	self.iCount = 1
	self.data = make(map[int]*BaseInfo)
	self.fun = make(map[int]TickFun)
}

func (self *TickControl) Delete() {
	close(self.C)
	self.C = nil

	self.data = map[int]*BaseInfo{}
	self.fun = make(map[int]TickFun)
}

func CreateTick() *TickControl {
	p := &TickControl{}
	p.Init()
	return p
}
