package router

import (
	"github.com/gin-gonic/gin"
)

func StartRouter() error {
	r := gin.Default()

	initializeGets(r)

	return r.Run("localhost:8080")
}
