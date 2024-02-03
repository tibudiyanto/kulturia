package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"kulturia/db"
	"kulturia/views"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

// Move this somewhere else or as another executable
func prepareDatabase(dbx db.DBTX) {
	fmt.Println("PREPARING DB")
	ctx := context.Background()
	if _, err := dbx.ExecContext(ctx, ddl); err != nil {
		fmt.Println("err", err)
	}
}

func main() {
	router := gin.Default()
	dbx, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		fmt.Println("OL", err)
	}

	go prepareDatabase((dbx))

	queries := db.New(dbx)

	// Home page
	router.GET("/", func(ctx *gin.Context) {
		entries, err := queries.GetEntries(ctx)
		page := views.Index(entries)
		template := views.Template("Home", page)
		if err != nil {
			fmt.Println("ERR FETCHING")
		}

		template.Render(ctx, ctx.Writer)
	})

	router.GET("/add", func(ctx *gin.Context) {
		page := views.Add("")
		template := views.Template("Tambah", page)

		template.Render(ctx, ctx.Writer)
	})

	type CreateForm struct {
		Name   string `form:"name"`
		Origin string `form:"origin"`
		Desc   string `form:"desc"`
	}

	router.POST("/add", func(ctx *gin.Context) {
		var entry CreateForm
		if err := ctx.ShouldBind(&entry); err != nil {
			fmt.Println("ERR", err)
			return
		}
		_, err := queries.CreateEntry(ctx, db.CreateEntryParams{
			Name:   entry.Name,
			Origin: entry.Origin,
			Desc:   entry.Desc,
		})

		if err != nil {
			page := views.Add(err.Error())
			page.Render(ctx, ctx.Writer)
			return
		}
		ctx.Header("HX-Redirect", "/")
		page := views.Add("Created")
		page.Render(ctx, ctx.Writer)
		return

	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
