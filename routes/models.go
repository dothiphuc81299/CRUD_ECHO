package routes

import (
	"echo/controllers"

	"github.com/labstack/echo"
)

//TodoRoute to...
func TodoRoute(g *echo.Group) {
	g.POST("", controllers.CreateModel)
	g.GET("", controllers.GetAllModel)
	g.GET("/:id", controllers.GetModelByID)
	g.DELETE("/:id", controllers.DeleteModel)
	g.PUT("/:id", controllers.UpdateModel)
	g.PATCH("/:id", controllers.CompletedModel)
}
