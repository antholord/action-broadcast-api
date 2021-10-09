package main

import (
	"fmt"
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

	var manager = NewManager()

	router.GET("/create", func(c *gin.Context) {
		fmt.Println("Create route called")

		clientName := c.DefaultQuery("user", "")
		if (clientName == "") {
			c.String(http.StatusBadRequest, "Cannot join session, user or sessionId missing")
			return
		}
		manager.HandleCreate(clientName, c)
	})
	router.GET("/join/:sessionId", func(c *gin.Context) {
		
		sessionId := c.Param("sessionId")
		if sessionId == "" {
			sessionId = "__default"
			// log.Fatalln("Missing SessionId to join session.")
			// return
		}
		clientName := c.DefaultQuery("user", "")
		if (clientName == "") {
			c.String(http.StatusBadRequest, "Cannot join session, user or sessionId missing")
			return
		}
		manager.HandleJoin(sessionId, clientName, c)
	})

	router.GET("session")

	router.Run(":" + port)
}
