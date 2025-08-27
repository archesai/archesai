package main

import (
	"embed"
	"net/http"

	"github.com/archesai/archesai/gen/api"
	"github.com/labstack/echo/v4"
)

//go:embed index.html
var swaggerUI embed.FS

func main() {

	e := echo.New()

	swagger, err := api.GetSwagger()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// serve swagger spec as JSON
	e.GET("/openapi.yaml", func(c echo.Context) error {
		data, err := swagger.MarshalJSON()
		if err != nil {
			return err
		}
		return c.JSONBlob(http.StatusOK, data)
	})

	// serve swagger docs
	e.GET("/swagger", echo.WrapHandler(http.StripPrefix("/swagger", http.FileServer(http.FS(swaggerUI)))))

	e.Start("0.0.0.0:3001")
}
