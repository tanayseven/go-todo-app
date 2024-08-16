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

const (
	TODO = 0
	EDIT = 1
	DONE = 2
)

type ListItemTable struct {
	gorm.Model
	ID    int
	Text  string
	State int
}

func (ListItemTable) TableName() string {
	return "list_item"
}

type ListItemViewModel struct {
	ID    int
	Text  string
	State string
}

func (l ListItemTable) ToListItemViewModel() ListItemViewModel {
	stringState := ""
	switch l.State {
	case TODO:
		stringState = "TODO"
		break
	case EDIT:
		stringState = "EDIT"
		break
	case DONE:
		stringState = "DONE"
	}
	return ListItemViewModel{
		ID:    l.ID,
		Text:  l.Text,
		State: stringState,
	}
}

func ToListItemViewModel(listItems []ListItemTable) []ListItemViewModel {
	var listItemsViewModel []ListItemViewModel
	for _, listItem := range listItems {
		listItemsViewModel = append(listItemsViewModel, listItem.ToListItemViewModel())
	}
	return listItemsViewModel
}

type TodoBaseViewModel struct {
	Name      string              `json:"message"`
	ListItems []ListItemViewModel `json:"products"`
}

type TodoItemViewModel struct {
	ID    int
	Text  string
	State string
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
		templates: template.Must(template.ParseGlob("templates/*.gohtml")),
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

	e.GET("/", func(c echo.Context) error {
		var listItems []ListItemTable
		db.Find(&listItems)
		return c.Render(
			http.StatusOK,
			"todo-base.gohtml",
			TodoBaseViewModel{
				Name:      username,
				ListItems: ToListItemViewModel(listItems),
			})
	})

	e.POST("/todo/add", func(c echo.Context) error {
		result := &ListItemTable{
			Text:  c.FormValue("text"),
			State: TODO,
		}
		db.Create(&result)
		return c.Render(http.StatusOK, "todo-item.gohtml", result.ToListItemViewModel())
	})

	e.DELETE("/todo/:id", func(c echo.Context) error {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Delete(&listItem)
		return c.NoContent(http.StatusOK)
	})

	e.PATCH("/todo/:id/edit", func(c echo.Context) error {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("State", EDIT)
		return c.Render(http.StatusOK, "todo-item.gohtml", listItem.ToListItemViewModel())
	})

	e.PATCH("/todo/:id/edit/save", func(c echo.Context) error {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("Text", c.FormValue("text"))
		db.Model(&listItem).Update("State", TODO)
		return c.Render(http.StatusOK, "todo-item.gohtml", listItem.ToListItemViewModel())
	})

	e.PATCH("/todo/:id/edit/cancel", func(c echo.Context) error {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("State", TODO)
		return c.Render(http.StatusOK, "todo-item.gohtml", listItem.ToListItemViewModel())
	})

	e.PATCH("/todo/:id/done", func(c echo.Context) error {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("State", DONE)
		return c.Render(http.StatusOK, "todo-item.gohtml", listItem.ToListItemViewModel())
	})

	e.PATCH("/todo/:id/undo", func(c echo.Context) error {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("State", TODO)
		return c.Render(http.StatusOK, "todo-item.gohtml", listItem.ToListItemViewModel())
	})
	e.Logger.Fatal(e.Start(":8000"))
}
