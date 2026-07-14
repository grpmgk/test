package calculator

import (
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/v4"
	"calculator\calculations"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculations)

	e.Start("localhost:8080")
}
