package main

import (
	"encoding/json"
	_ "github.com/GoAdminGroup/go-admin/adapter/echo" // Import the adapter, it must be imported. If it is not imported, you need to define it yourself.
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite" // Import the sql driver
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/tests/tables"
	_ "github.com/GoAdminGroup/themes/adminlte" // Import the theme
	"html/template"
	"io"

	"net/http"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Message struct {
	Message string `json:"message"`
}

func main() {
	cfg := config.Config{
		Databases: config.DatabaseList{
			"default": {
				File:   "main.db",
				Driver: config.DriverSqlite,
			},
		},
		UrlPrefix: "admin", // The url prefix of the website.
		// Store must be set and guaranteed to have write access, otherwise new administrator users cannot be added.
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language: language.EN,
	}

	e := echo.New()
	eng := engine.Default()
	eng.AddConfig(&cfg).
		AddGenerators(tables.Generators)

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.gohtml")),
	}

	// Set the renderer
	e.Renderer = renderer

	_ = eng.Use(e)

	newMap := make(map[string]interface{})
	data, _ := json.Marshal(Message{Message: "Hello, World!"})
	_ = json.Unmarshal(data, &newMap)
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "base.gohtml", newMap)
	})

	e.Logger.Fatal(e.Start(":8081"))
}
