package leader

import (
	"net"
	"context"

	"xqdfs/constant"
	"xqdfs/utils/plugin"
	"xqdfs/utils/helper"

	"github.com/Jeffail/gabs"
)

const(
	CmdLeader 	= "/strategy/leader/get"
)

var(
	leaderSelect *LeaderSelect
)

func ServiceLeaderSelectSetup(leader *LeaderSelect) {
	leaderSelect=leader
	plugin.PluginAddService(CmdLeader,ServiceLeaderGet)
}

func ServiceLeaderGet(ctx context.Context,inv *plugin.Invocation) interface{}{
	var ip string
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			if ipnet, flag := address.(*net.IPNet); flag && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip=ipnet.IP.String()
					break
				}
			}
		}
	}

	json:=gabs.New()
	json.Set(leaderSelect.leaderId,"leaderId")
	json.Set(leaderSelect.meId,"meId")
	json.Set(ip,"addr")
	return helper.ResultBuildWithBody(constant.Success,json)
}