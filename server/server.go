package server

import (
	"github.com/zwpaper/godback/utils"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
)

var HTTPServer *gin.Engine
var roomHub map[string]*Game
var logger *logs.BeeLogger

func init() {
	logger = utils.Log
	roomHub = map[string]*Game{}
	HTTPServer = gin.Default()
	setRoute(HTTPServer)
}
