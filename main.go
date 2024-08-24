package main

import (
	rice "github.com/GeertJohan/go.rice"
	_ "github.com/GoAdminGroup/go-admin/adapter/echo" // Import the adapter, it must be imported. If it is not imported, you need to define it yourself.
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite" // Import the sql driver
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/tests/tables"
	_ "github.com/GoAdminGroup/themes/adminlte" // Import the theme
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"io"
	"net/http"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	eng.AddConfig(&cfg).
		AddGenerators(tables.Generators)
	// Create a new rice.Box for the templates
	templateBox, err := rice.FindBox("templates")
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Load all the templates from rice-box.go
	t := &Template{
		templates: template.Must(template.New("main").Parse(templateBox.MustString("base.gohtml"))),
	}
	staticBox, err := rice.FindBox("static")
	if err != nil {
		e.Logger.Fatal(err)
	}
	username := "john.doe"
	e.Renderer = t
	e.Debug = true
	e.Static("/static", staticBox.Name())
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
