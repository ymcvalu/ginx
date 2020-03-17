package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ymcvalu/ginx"
	"log"
)

func main() {
	app := gin.New()
	router := ginx.XRouter(app, nil)
	router.Any("/greeting", func(who *struct {
		Name string `json:"name" form:"name" binding:"required"`
	}) string {
		return fmt.Sprintf("Hi, %s!", who.Name)
	})

	log.Fatal(app.Run(":8090"))
}
