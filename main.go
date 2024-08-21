package main

import (
	_ "github.com/GoAdminGroup/go-admin/adapter/echo" // Import the adapter, it must be imported. If it is not imported, you need to define it yourself.
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite" // Import the sql driver
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/tests/tables"
	_ "github.com/GoAdminGroup/themes/adminlte" // Import the theme
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func registerRoot(e *echo.Echo, db *gorm.DB, username string) *echo.Route {
	return e.GET("/", func(c echo.Context) error {
		var listItems []ListItemTable
		db.Find(&listItems)
		return c.Render(
			http.StatusOK,
			"base.gohtml",
			TodoBaseViewModel{
				Name:      username,
				ListItems: ToListItemViewModel(listItems),
			})
	})
}

func main() {
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	cfg := config.Config{
		Databases: config.DatabaseList{
			"default": {
				File:   "main.db",
				Driver: config.DriverSqlite,
			},
		},
		UrlPrefix: "admin",
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
		templates: template.Must(template.ParseGlob("templates/**/*.gohtml")),
	}
	username := "john.doe"
	e.Renderer = renderer
	e.Debug = true
	e.Static("/static", "static")
	e.Use(middleware.Logger())
	_ = eng.Use(e)

	err = db.AutoMigrate(ListItemTable{})
	if err != nil {
		return
	}

	registerRoot(e, db, username)
	registerTodos(e, db)

	e.Logger.Fatal(e.Start(":8000"))
}
