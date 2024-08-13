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
	Message  string        `json:"message"`
	Products []ProductData `json:"products"`
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
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
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
	var products []Product
	db.Find(&products)

	e.GET("/", func(c echo.Context) error {
		return c.Render(
			http.StatusOK,
			"base.gohtml",
			Data{
				Message:  "Hello, World!",
				Products: ToProductDataList(products),
			})
	})

	e.Logger.Fatal(e.Start(":8000"))
}
