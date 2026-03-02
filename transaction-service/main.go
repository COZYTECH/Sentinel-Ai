package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	InitDB()
	r := gin.Default()

	r.POST("/transactions", CreateTransaction)

	r.Run(":8081")
}