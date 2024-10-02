package router

import "github.com/gin-gonic/gin"

func GetInit(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/api/articles", func(c *gin.Context) {
		c.HTML(200, "articles.html", gin.H{
			"title": "Articles",
		})
	})
}
