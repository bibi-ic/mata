package main

import (
	"embed"
	"io/fs"
	"log"

	"github.com/bibi-ic/mata/config"
	"github.com/bibi-ic/mata/internal/api"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:embed all:swagger-ui
var templateFS embed.FS

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	staticFile, err := fs.Sub(templateFS, "swagger-ui")
	if err != nil {
		log.Fatal("cannot load static swagger: ", err)
	}
	c.InsertFS(staticFile)

	s := api.NewServer(c)
	err = s.Start()
	if err != nil {
		log.Fatal("can not run server: ", err)
	}
}
