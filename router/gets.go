package router

import (
	templates "github.com/Bolado/ai-tracker/website/templates"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func initializeGets(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		templ.Handler(templates.Index()).ServeHTTP(c.Writer, c.Request)
	})

	r.Static("/static", "./website/static")
}
