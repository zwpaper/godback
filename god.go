package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/zwpaper/godinwerewolves/server"
	"github.com/zwpaper/godinwerewolves/store"
	"github.com/zwpaper/godinwerewolves/utils"
)

func init() {
	logs.SetLevel(utils.LogLevel[utils.Conf.LogConf.Level])
	logs.EnableFuncCallDepth(true)
}

func main() {
	logs.Info("Welcome to God of werewolves")
	logs.Debug("Config: %v", utils.Conf)
	err := utils.SetBell(utils.Conf.AlarmConf.Url, utils.Conf.AlarmConf.Module,
		[]string{},
		utils.Conf.AlarmConf.Timeout, 1, utils.Conf.AlarmConf.IsOn)
	if err != nil {
		logs.Emergency("Can not set Alarm: %v", err)
	}

	store.Init(utils.Conf.EtcdConf.BindAddrs, utils.Conf.EtcdConf.Prefix)
	god := server.HTTPServer
	logs.Emergency(god.Run(":" + utils.Conf.HTTPConf.Port))
	//logs.Emergency(r.RunTLS("0.0.0.0:8080", "server.pem", "server.key")) // listen and serve on 0.0.0.0:8080
}
