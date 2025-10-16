package main

import (
	"supernova/authService/auth"

	"github.com/gin-gonic/gin"
)

func main() {


	router := gin.Default()

	auth.SetupAuthApp(router)

	router.Run(":8081")
}