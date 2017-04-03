package server

import "github.com/gin-gonic/gin"

var HTTPServer *gin.Engine
var roomHub map[string]*Hub

func init() {
	roomHub = map[string]*Hub{}
	HTTPServer = gin.Default()
	setRoute(HTTPServer)
}
