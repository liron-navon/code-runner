package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"
	"github.com/liron-navon/code-runner/lib/handlers"
)

func main() {
	ginEngine := gin.Default()

	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	ginEngine.POST("/exec", handlers.HandleExec)
	ginEngine.Run() // listen and serve on 0.0.0.0:8080
}
