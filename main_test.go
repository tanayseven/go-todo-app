package main

import (
	"github.com/andybalholm/cascadia"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func Query(n *html.Node, query string) *html.Node {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return &html.Node{}
	}
	return cascadia.Query(n, sel)
}

func QueryAll(n *html.Node, query string) []*html.Node {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return []*html.Node{}
	}
	return cascadia.QueryAll(n, sel)
}

func GetAttributeValue(input *html.Node, key string) string {
	value := ""
	for _, a := range input.Attr {
		if a.Key == key {
			value = a.Val
		}
	}
	return value
}

func setupTest(tb testing.TB) func(tb testing.TB) {
	_ = os.Remove("test.db")

	return func(tb testing.TB) {
		_ = os.Remove("test.db")
	}
}

func TestListIsInitiallyEmpty(t *testing.T) {
	defer (setupTest(t))(t)
	db := setupDatabase("test.db")
	router := setupServer(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	doc, err := html.Parse(strings.NewReader(w.Body.String()))
	assert.Nil(t, err)
	assert.NotNil(t, doc)
	ul := Query(doc, "ul")
	lis := QueryAll(ul, "li")
	assert.Equal(t, 0, len(lis))
}

func TestAddItem(t *testing.T) {
	defer (setupTest(t))(t)
	db := setupDatabase("test.db")
	router := setupServer(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/todo/add", strings.NewReader("text=Buy Milk"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	doc, err := html.Parse(strings.NewReader(w.Body.String()))
	assert.Nil(t, err)
	assert.NotNil(t, doc)
	input := Query(doc, "li > input")
	assert.Equal(t, "Buy Milk", GetAttributeValue(input, "value"))
}

func TestDeleteItem(t *testing.T) {
	defer (setupTest(t))(t)
	db := setupDatabase("test.db")
	router := setupServer(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/todo/add", strings.NewReader("text=Buy Milk"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/todo/1", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	doc, err := html.Parse(strings.NewReader(w.Body.String()))
	assert.Nil(t, err)
	assert.NotNil(t, doc)
	ul := Query(doc, "ul")
	lis := QueryAll(ul, "li")
	assert.Equal(t, 0, len(lis))
}
