package main

import (
	rice "github.com/GeertJohan/go.rice"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
	_ "github.com/GoAdminGroup/themes/adminlte"
	"github.com/foolin/gin-template/supports/gorice"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func registerRoot(e *gin.Engine, db *gorm.DB, username string) {
	e.GET("/", func(c *gin.Context) {
		var listItems []ListItemTable
		db.Find(&listItems)
		c.HTML(http.StatusOK, "base.gohtml", TodoBaseViewModel{
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
	gingine := gin.Default()

	username := "john.doe"
	gingine.HTMLRender = gorice.New(rice.MustFindBox("templates"))
	gingine.LoadHTMLGlob("templates/**/*.gohtml")
	staticBox := rice.MustFindBox("static")
	gingine.StaticFS("/static", staticBox.HTTPBox())

	_ = db.AutoMigrate(ListItemTable{})

	registerRoot(gingine, db, username)
	registerTodos(gingine, db)

	_ = gingine.Run(":9033")
}
