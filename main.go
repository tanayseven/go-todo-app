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
	State int
}

func (l ListItemTable) ToListItemViewModel() ListItemViewModel {
	return ListItemViewModel{
		ID:    l.ID,
		Text:  l.Text,
		State: l.State,
	}
}

func ToListItemViewModelList(listItems []ListItemTable) []ListItemViewModel {
	var listItemsViewModel []ListItemViewModel
	for _, listItem := range listItems {
		listItemsViewModel = append(listItemsViewModel, listItem.ToListItemViewModel())
	}
	return listItemsViewModel
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type ProductData struct {
	Code  string
	Price uint
}

func (p Product) ToProductData() ProductData {
	return ProductData{
		Code:  p.Code,
		Price: p.Price,
	}
}

func ToProductDataList(products []Product) []ProductData {
	var productsData []ProductData
	for _, product := range products {
		productsData = append(productsData, product.ToProductData())
	}
	return productsData
}

type Data struct {
	Message   string              `json:"message"`
	ListItems []ListItemViewModel `json:"products"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&Product{})
	if err != nil {
		return
	}
	db.Create(&Product{Code: "D42", Price: 100})
	var product Product
	db.First(&product, 1)
	db.First(&product, "code = ?", "D42")
	db.Model(&product).Update("Price", 200)
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"})
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
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
	e.Renderer = renderer
	e.Debug = true
	e.Use(middleware.Logger())
	_ = eng.Use(e)
	var listItems []ListItemTable
	db.Find(&listItems)

	e.GET("/", func(c echo.Context) error {
		return c.Render(
			http.StatusOK,
			"base.gohtml",
			Data{
				Message:   "Hello, World!",
				ListItems: ToListItemViewModelList(listItems),
			})
	})

	//e.DELETE("/products/:id", func(c echo.Context) error {
	//	// TODO: delete an entry
	//}
	//
	//e.POST("/products", func(c echo.Context) error {
	//	// TODO: create a new entry
	//}
	//
	//e.PATCH("/products/:id/edit", func(c echo.Context) error {
	//	// TODO: edit the text of an entry
	//}
	//
	//e.PATCH("/products/:id/done", func(c echo.Context) error {
	//	// TODO: mark the text of an entry as done or not done
	//}
	e.Logger.Fatal(e.Start(":8000"))
}
