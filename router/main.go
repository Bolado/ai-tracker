package router

import (
	"github.com/gin-gonic/gin"
)

func StartRouter() error {
	r := gin.Default()

	initializeGets(r)

	return r.Run(":8080")
}
