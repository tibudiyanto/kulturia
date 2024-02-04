package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"kulturia/config"
	"kulturia/db"
	"kulturia/views"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type R2Client struct {
	s3     *s3.Client
	config config.Config
}

func (r *R2Client) GetBucketLocation(id int64) string {
	return fmt.Sprintf(r.config.R2PublicURL, id)
}

func getR2(r2Config config.Config) *R2Client {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r2Config.R2AccountId),
		}, nil
	})

	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithEndpointResolverWithOptions(r2Resolver),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(r2Config.R2AccessKeyId, r2Config.R2AccessKeySecret, "")),
		awsConfig.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)
	return &R2Client{
		s3:     client,
		config: r2Config,
	}
}

func main() {
	config, err := config.GetConfig()

	if err != nil {
		panic("Need .env file")
	}

	router := gin.Default()
	dbx, err := sql.Open("sqlite3", "db.sqlite3")

	if err != nil {
		panic(err)
	}

	r2Client := getR2(config)

	queries := db.New(dbx)

	// Home page
	router.GET("/", func(ctx *gin.Context) {
		entries, err := queries.GetEntries(ctx)
		page := views.Index(entries)
		template := views.Template("Home", page)
		if err != nil {
			fmt.Println("ERR FETCHING", err)
		}

		template.Render(ctx, ctx.Writer)
	})

	type MemeURI struct {
		ID int64 `uri:"id" binding:"required"`
	}

	router.GET("/:id", func(ctx *gin.Context) {
		var uri MemeURI
		if err := ctx.ShouldBindUri(&uri); err != nil {
			fmt.Println("OI", err)
		}

		entry, err := queries.GetEntry(ctx, uri.ID)

		if err != nil {
			// 404
			ctx.Redirect(404, "404")
		}

		page := views.Show(entry)
		template := views.Template("Tambah", page)

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

		createdEntry, err := queries.CreateEntry(ctx, db.CreateEntryParams{
			Name:   entry.Name,
			Origin: entry.Origin,
			Desc:   entry.Desc,
		})

		if err != nil {
			panic(err)
		}

		file, err := ctx.FormFile("asset")
		if err != nil {
			panic(err)
		}
		fileContent, err := file.Open()
		if err != nil {
			panic(err)
		}
		defer fileContent.Close()

		_, err = r2Client.s3.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(config.R2BucketName),
			Key:    aws.String(strconv.FormatInt(createdEntry.ID, 10)),
			Body:   fileContent,
		})

		if err != nil {
			fmt.Println("OI", err)
		}

		_, err = queries.CreateAsset(ctx, db.CreateAssetParams{
			EntryID:  sql.NullInt64{Int64: createdEntry.ID, Valid: true},
			Location: sql.NullString{String: r2Client.GetBucketLocation(createdEntry.ID), Valid: true},
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
