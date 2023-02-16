package main

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/gorm"
	"hamta-sharif/models"
	"hamta-sharif/repository"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	db := repository.Initialize()
	h := newPostgresRepo(db)
	setUpRouts(h)
}

type postgresRepo struct {
	DB *gorm.DB
}

func newPostgresRepo(DB *gorm.DB) *postgresRepo {
	return &postgresRepo{DB: DB}
}

func setUpRouts(h *postgresRepo) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "static")

	e.GET("/create-file/:id", h.DownloadFile)
	e.POST("/create-file", h.uploadFile)
	e.Logger.Fatal(e.Start(":8000"))
}

func (p *postgresRepo) uploadFile(c echo.Context) error {
	var req models.CreateFileReq
	err := c.Bind(&req)
	if err != nil {
		return echo.ErrBadRequest
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			return
		}
	}(src)

	// Destination
	dst, err := os.Create("static/" + req.Name)
	if err != nil {
		return err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			return
		}
	}(dst)

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	newFile, err := p.createFile(req)
	return c.String(http.StatusCreated, strconv.Itoa(newFile.Id))
}

func (p *postgresRepo) createFile(req models.CreateFileReq) (*models.File, error) {
	newFile := &models.File{
		Name:       req.Name,
		Version:    req.Version,
		DateCreate: time.Now(),
	}

	err := p.DB.Create(newFile).Error
	if err != nil {
		return nil, echo.ErrBadRequest
	}
	return newFile, nil
}

func (p *postgresRepo) DownloadFile(c echo.Context) error {
	id := c.Param("id")
	var file *models.File
	err := p.DB.Where("id= ?", id).Find(&file)
	if err != nil {
		return errors.New("this file is not exist")
	}
	return c.Redirect(http.StatusTemporaryRedirect, file.Name)
}
