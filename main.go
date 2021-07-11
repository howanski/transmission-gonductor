package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	router = gin.Default()

	router.LoadHTMLGlob("front-end/html/*")

	router.Static("/assets/images", "./front-end/images")
	router.Static("/assets/js", "./front-end/js")

	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "GOnductor",
			},
		)
	})

	router.GET("/gonductor-stats", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"connectionStatus": "Ruined",
				"lastPing":         "Never",
			},
		)
	})

	router.Run()
}
