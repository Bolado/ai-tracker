package router

import (
	"strconv"

	templates "github.com/Bolado/ai-tracker/website/templates"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func initializeGets(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		//check for page parameter
		pageQuery := c.Query("page")

		page := 0
		if pageQuery != "" {
			// check if page is a number and parse it
			if p, err := strconv.Atoi(pageQuery); err != nil {

				// if it's not a number, set page to 0
				page = 0

			} else {

				// if page is negative, set it to 0
				if p < 0 {
					page = 0
				}

				// if page is greater than 0, set it to 0
				if p > 0 {
					page = p
				}

				if p > templates.GetNumberOfPages()-1 {
					page = templates.GetNumberOfPages() - 1
				}
			}
		}

		templ.Handler(templates.Index(page)).ServeHTTP(c.Writer, c.Request)
	})

	r.Static("/static", "./website/static")
}
