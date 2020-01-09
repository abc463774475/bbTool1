package n_rpc

import (
	"bbTool/n_log"
	"net"
	"net/rpc"
	"sync"
	"bbTool/n_json"
)

// 一个数据的连接
type oneConnectMgr struct {
	free  map[string]*poolBaseData // 空闲
	using map[string]*poolBaseData // 使用中

	addr string // 地址族
}

func CreateOneConnectMgr(addr string) *oneConnectMgr {
	p := &oneConnectMgr{}
	p.addr = addr
	p.free = make(map[string]*poolBaseData)
	p.using = make(map[string]*poolBaseData)

	return p
}

func (self *oneConnectMgr) GetFree() *poolBaseData {
	if len(self.free) == 0 {
		if len(self.using) >= 20 {
			// 连接数过多咋办 ？？
			// 咋个办 呢 ？
			return nil
		}
		p := CreateBaseData(self.addr)
		self.free[p.index] = p
	}

	n_log.OnlyFile("free  %v  using %v", len(self.free), len(self.using))
	for k, v := range self.free {
		self.using[k] = v
		delete(self.free, k)
		return v
	}
	return nil
}

type RemoteFamilyMgr struct {
	all map[string]*oneConnectMgr // 所有数据
	l   sync.Mutex
}

func CreateMgr() *RemoteFamilyMgr {
	p := &RemoteFamilyMgr{}
	p.all = make(map[string]*oneConnectMgr)

	return p
}

func (self *RemoteFamilyMgr) GetRemoteClient(addr string) *poolBaseData {
	self.l.Lock()
	defer self.l.Unlock()

	m, ok := self.all[addr]
	if !ok {
		p := CreateOneConnectMgr(addr)
		self.all[addr] = p
	}

	m = self.all[addr]
	return m.GetFree()
}

// 删除原有数据
func (self *RemoteFamilyMgr) AddFreeClient(addr, index string, p *poolBaseData) {
	self.l.Lock()
	defer self.l.Unlock()

	m := self.all[addr]
	delete(m.using, index)
	m.free[index] = p
}

func CallRemoteInterface(addr string, fName string, i interface{}, reply *[]byte) error{
	data ,_ := n_json.Marshal(i)
	return CallRemotebase(addr,fName,data,reply)
}

// db sala gs cs
func CallRemotebase(addr string, fName string, args []byte, reply *[]byte) error {
	var pClient *poolBaseData

	for {
		pClient = G_RemotePool.GetRemoteClient(addr)
		if pClient != nil {
			break
		}
	}
	// 连接不上

	if pClient.Client == nil { // 初始化的时候   该怎样弄呢 ？？  哥哥
		address, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			panic(err)
		}
		conn, e := net.DialTCP("tcp", nil, address)
		if e != nil {
			return e
		}
		n_log.Erro("err  %v  conn %v",e,conn)
		pClient.Client = rpc.NewClient(conn)
	}

	err := pClient.Call(fName, args, reply)
	if err != nil {
		n_log.Erro("rpc erro  %v",err)
	}

	return err
}

func Print_Info() {
	G_RemotePool.l.Lock()
	defer G_RemotePool.l.Unlock()

	for _, v := range G_RemotePool.all {
		n_log.Debug("total free len %v  total using %v", len(v.free), len(v.using))
		n_log.Erro("total free %v  using %v", v.free, v.using)
	}
}
