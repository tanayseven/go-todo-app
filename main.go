package main

import (
	"fmt"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
	_ "github.com/GoAdminGroup/themes/adminlte"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"os"
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
	db := setupDatabase("main.db")
	gingine := setupServer(db)
	portNumber := os.Getenv("PORT")
	if portNumber == "" {
		portNumber = "9033"
	}

	_ = gingine.Run(fmt.Sprintf(":%s", portNumber))
}

func setupDatabase(databasePath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func setupServer(db *gorm.DB) *gin.Engine {
	gingine := gin.Default()

	username := "john.doe"
	gingine.LoadHTMLGlob("templates/**/*.gohtml")
	gingine.StaticFS("/static", http.Dir("static"))

	_ = db.AutoMigrate(ListItemTable{})

	registerRoot(gingine, db, username)
	registerTodos(gingine, db)
	return gingine
}
