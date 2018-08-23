package main

import (
	"fmt"
	"log"
	"path"
	"runtime"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/henrocdotnet/golang-playground/internal/note-monger/handler"
)



type CustomContext struct {
	echo.Context
}

func main() {
	log.Println("BEGIN: main")

	e := echo.New()

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c}
			return h(cc)
		}
	})



	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(filename)
	fmt.Printf("filename: %s, dir: %s\n", filename, dir)
	path := path.Join(path.Dir(filename), "../../web/note-monger/static")
	// e.Use(middleware.Static("../../web/note-monger/static"))
	e.Use(middleware.Static(path))

	e.GET("/", handler.Index)

	/*
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ROOT HANDLER")
	})
	*/
	e.Logger.Fatal(e.Start(":8080"))
}

