package main

import (
	"kulturia/views"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		page := views.HelloContent("OLA")
		template := views.Template("Home", page)
		template.Render(ctx, ctx.Writer)
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
