package utils

import "github.com/astaxie/beego/logs"

var Log *logs.BeeLogger
var log *logs.BeeLogger

func init() {
	log = logs.NewLogger()
	err := log.SetLogger(logs.AdapterConsole)
	if err != nil {
		panic("Can not set logger")
	}
	log.EnableFuncCallDepth(true)

	Log = log
}
