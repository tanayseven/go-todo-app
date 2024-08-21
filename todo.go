package main

import (
	"github.com/labstack/echo/v4"
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

func registerTodos(e *echo.Echo, db *gorm.DB) {
	e.POST("/todo/add", func(c echo.Context) error {
		result := &ListItemTable{
			Text:  c.FormValue("text"),
			State: TODO,
		}
		db.Create(&result)
		return c.Render(http.StatusOK, "item.gohtml", result.ToListItemViewModel())
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
		return c.Render(http.StatusOK, "item.gohtml", listItem.ToListItemViewModel())
	})

	e.PATCH("/todo/:id/edit/save", func(c echo.Context) error {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("Text", c.FormValue("text"))
		db.Model(&listItem).Update("State", TODO)
		return c.Render(http.StatusOK, "item.gohtml", listItem.ToListItemViewModel())
	})

	e.PATCH("/todo/:id/edit/cancel", func(c echo.Context) error {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("State", TODO)
		return c.Render(http.StatusOK, "item.gohtml", listItem.ToListItemViewModel())
	})

	e.PATCH("/todo/:id/done", func(c echo.Context) error {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("State", DONE)
		return c.Render(http.StatusOK, "item.gohtml", listItem.ToListItemViewModel())
	})

	e.PATCH("/todo/:id/undo", func(c echo.Context) error {
		var listItem ListItemTable
		db.First(&listItem, c.Param("id"))
		db.Model(&listItem).Update("State", TODO)
		return c.Render(http.StatusOK, "item.gohtml", listItem.ToListItemViewModel())
	})
}
