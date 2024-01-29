package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bibi-ic/mata/internal/models"
)

func RandMeta(url string) models.Meta {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal("cannot get working dir: ", err)
	}
	parent := filepath.Dir(path)
	root := filepath.Dir(parent)

	data, err := os.ReadFile(root + "/test/testdata/testMeta.data")
	if err != nil {
		log.Fatal("cannot get sample metadata: ", err)
	}

	var metas []models.Meta
	parseErr := json.Unmarshal(data, &metas)
	if parseErr != nil {
		fmt.Println(metas)
		if e, ok := parseErr.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
			log.Println(string(data[0:e.Offset]))
		}

		log.Fatal("cannot parse json data: ", parseErr)
	}

	for _, meta := range metas {
		if meta.URL == url {
			fmt.Println(url)
			return meta
		}
	}

	log.Fatal("cannot find sample metadata match url given")
	return models.Meta{}
}
