package n_rpc

import "bbTool/n_log"

var (
	G_RemotePool *RemoteFamilyMgr
)

func init() {
	G_RemotePool = CreateMgr()
	n_log.Debug("rpc  init ok")
}
