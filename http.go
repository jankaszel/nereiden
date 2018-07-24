package main

import (
	"github.com/gin-gonic/gin"
)

func createSecuredGroup(router *gin.Engine, authUser string, authPassword string) *gin.RouterGroup {
	if authUser == "" || authPassword == "" {
		return router.Group("/")
	}

	accounts := gin.Accounts{}
	accounts[authUser] = authPassword

	return router.Group("/", gin.BasicAuth(accounts))
}
