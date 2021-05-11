package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	. "github.com/heroku/action-broadcast-api/src/websockets"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/russross/blackfriday"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, string(blackfriday.Run([]byte("**hi!**"))))
	})

	var hub = NewHub()
	go hub.Run()

	router.GET("/ws", func(c *gin.Context) {
		ServeWs(hub, c.Writer, c.Request)
	})

	router.Run(":" + port)
}
