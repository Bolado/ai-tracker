package router

import (
	"github.com/gin-gonic/gin"
)

func StartRouter() {
	r := gin.Default()

	GetInit(r)

	r.Run(":8080")
}
