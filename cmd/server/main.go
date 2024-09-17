package main

import (
	"github.com/gin-gonic/gin"
	"github.com/salberternst/portal/pkg/routes"
)

func main() {
	r := gin.Default()

	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/manifest.json", "./public/manifest.json")
	r.StaticFile("/cropped-android-chrome-512x512-1-1-192x192.png", "./public/cropped-android-chrome-512x512-1-1-192x192.png")
	r.StaticFile("/config.js", "./public/config.js")
	r.Static("/assets", "./public/assets")

	routes.AddRoutes(r)

	r.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	r.Run()
}
