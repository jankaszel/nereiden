package main

import (
	"net/http"

	recreate "github.com/fallafeljan/docker-recreate/lib"
	"github.com/gin-gonic/gin"
)

func perform(c *gin.Context) {
	val, exists := c.Get("tokenConf")
	tokenConf, ok := val.(TokenConf)

	if !exists || !ok || tokenConf.ContainerID == "" || tokenConf.ImageTag == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal server error"})

		return
	}

	options := recreate.Options{
		PullImage:       true,
		DeleteContainer: true,
		Registries:      []recreate.RegistryConf{}}

	recreation, err := recreate.Recreate(
		"unix:///var/run/docker.sock",
		tokenConf.ContainerID,
		tokenConf.ImageTag,
		&options)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal server error",
			"error":  err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "ok",
		"fromContainer": recreation.PreviousContainerID[:8],
		"toContainer":   recreation.NewContainerID[:8]})
}
