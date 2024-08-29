package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

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

func (l ListItemTable) ToListItemJsonModel() ListItemJsonModel {
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
	return ListItemJsonModel{
		ID:    int(l.ID),
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

func registerTodos(e *gin.Engine, db *gorm.DB) {
	e.POST("/todo/add", func(c *gin.Context) {
		err := c.Request.ParseForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // TODO change this to an error html
		}
		if c.PostForm("text") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "text is empty"})
			return
		}
		result := &ListItemTable{
			Text:  c.PostForm("text"),
			State: TODO,
		}
		db.Create(&result)
		c.HTML(http.StatusOK, "item.gohtml", result.ToListItemViewModel())
	})

	e.DELETE("/todo/:id", func(c *gin.Context) {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Delete(&listItem)
		c.Data(http.StatusOK, gin.MIMEHTML, nil)
	})

	e.PATCH("/todo/:id/edit", func(c *gin.Context) {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("State", EDIT)
		c.HTML(http.StatusOK, "item.gohtml", listItem.ToListItemViewModel())
	})

	e.PATCH("/todo/:id/edit/save", func(c *gin.Context) {
		err := c.Request.ParseForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // TODO change this to an error html
		}
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("Text", c.PostForm("text"))
		db.Model(&listItem).Update("State", TODO)
		c.HTML(http.StatusOK, "item.gohtml", listItem.ToListItemViewModel())
	})

	e.PATCH("/todo/:id/edit/cancel", func(c *gin.Context) {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("State", TODO)
		c.HTML(http.StatusOK, "item.gohtml", listItem.ToListItemViewModel())
	})

	e.PATCH("/todo/:id/done", func(c *gin.Context) {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("State", DONE)
		c.HTML(http.StatusOK, "item.gohtml", listItem.ToListItemViewModel())
	})

	e.PATCH("/todo/:id/undo", func(c *gin.Context) {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("State", TODO)
		c.HTML(http.StatusOK, "item.gohtml", listItem.ToListItemViewModel())
	})
}

func registerTodoApi(e *gin.Engine, db *gorm.DB) {
	// @BasePath /api

	//TO-DO godoc
	// @Summary Get to-do items
	// @Description Get all to-do items
	// @Router /todo [get]
	// @Produce json
	// @Success 200 {object} []ListItemViewModel
	e.GET("/api/todo", func(c *gin.Context) {
		var listItems []ListItemTable
		db.Find(&listItems)
		c.JSON(http.StatusOK, ToListItemJSONModel(listItems))
	})

	e.POST("/api/todo", func(c *gin.Context) {
		var listItem ListItemTable
		err := c.BindJSON(&listItem)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&listItem)
		c.JSON(http.StatusOK, listItem.ToListItemViewModel())
	})

	e.PATCH("/api/todo/:id", func(c *gin.Context) {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		err := c.BindJSON(&listItem)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&listItem)
		c.JSON(http.StatusOK, listItem.ToListItemViewModel())
	})
}

type ListItemJsonModel struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	State string `json:"state"`
}

func ToListItemJSONModel(items []ListItemTable) []ListItemJsonModel {
	var listItems []ListItemJsonModel
	for _, item := range items {
		listItems = append(listItems, item.ToListItemJsonModel())
	}
	return listItems
}
