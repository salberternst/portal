package main

import (
	"github.com/gin-gonic/gin"
	"github.com/salberternst/fsn_landing_page/pkg/routes"
)

func main() {
	r := gin.Default()

	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/config.js", "./public/config.js")
	r.Static("/assets", "./public/assets")
	r.Static("/img", "./public/img")
	r.StaticFile("/404.html", "./public/404.html")

	routes.AddRoutes(r)

	r.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	r.Run()
}
