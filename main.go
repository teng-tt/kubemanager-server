package main

import (
	"kubmanager/global"
	"kubmanager/initiallize"
)

// 项目启动入口
func main() {
	r := initiallize.Routers()
	initiallize.Viper()
	initiallize.K8SWithDiscovery()
	panic(r.Run(global.CONF.System.Addr))
}
