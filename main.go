package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/p-mega/limit-flow/filter"
)

func main() {
	router := gin.Default()
	router.Use(filter.GinLimitFlow(2*1000, 10, 300*1000))
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	router.Run(":8000")
}
