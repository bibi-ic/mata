package main

import (
	"log"

	"github.com/bibi-ic/mata/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		log.Fatal("can not create server: ", err)
	}

	err = s.Start()
	if err != nil {
		log.Fatal("can not run server: ", err)
	}
}
