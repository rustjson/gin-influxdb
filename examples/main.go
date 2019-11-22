package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	gininflux "github.com/rustjson/gin-influxdb"
)

func main() {
	router := gin.New()

	// p := ginprometheus.NewPrometheus("gin")
	i := gininflux.New("http://localhost:8086", "test_jason_go", "gin", 2, map[string]string{})
	// p.Use(r)
	router.Use(i.HandlerFunc())
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world!")
	})
	v1 := router.Group("/v1")
	{
		v1.POST("/login", func(c *gin.Context) {
			c.JSON(200, "login")
		})
		v1.POST("/submit", func(c *gin.Context) {
			c.JSON(200, "submit")
		})
		v1.POST("/read", func(c *gin.Context) {
			c.JSON(200, "read")
		})
	}
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(200, "Hello !"+name)
	})

	fmt.Printf("r = %+v, router=%+v\n", router.Handlers[0], router.Routes()[0].Path)

	router.Run(":8000")
}
