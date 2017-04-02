package server

import (
	"plaid/utils"

	"github.com/gin-gonic/gin"
)

var HTTPServer *gin.Engine
var roomHub map[string]*Hub

func init() {
	log = utils.Log
	roomHub = make(map[string]*Hub)
	HTTPServer := gin.Default()
	setRoute(HTTPServer)
}
