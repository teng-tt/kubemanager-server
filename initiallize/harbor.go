package initiallize

import (
	"kubmanager/global"
	"kubmanager/plugins/harbor"
)

func InitHarborClient() {
	enable := global.CONF.System.Harbor.Enable
	cafile := global.CONF.System.Harbor.CacertPath
	scheme := global.CONF.System.Harbor.Scheme
	username := global.CONF.System.Harbor.Username
	password := global.CONF.System.Harbor.Password
	host := global.CONF.System.Harbor.Host
	initHarborClient, err := harbor.InitHarbor(scheme, host, username, password, cafile)
	if err != nil && enable {
		panic(err)
	}
	global.HarborClient = initHarborClient
}
