package main

import (
	"log"

	"github.com/bibi-ic/mata/config"
	"github.com/bibi-ic/mata/internal/api"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	s := api.NewServer(c)
	err = s.Start()
	if err != nil {
		log.Fatal("can not run server: ", err)
	}
}
