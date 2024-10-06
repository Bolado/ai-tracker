package router

import (
	templates "github.com/Bolado/aitracker/website/templates"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func GetInit(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		templ.Handler(templates.Index()).ServeHTTP(c.Writer, c.Request)
	})

	r.Static("/static", "./website/static")
}
