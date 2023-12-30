package randomutil

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

func RandURL() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal("cannot get working dir: ", err)
	}
	parent := filepath.Dir(path)
	root := filepath.Dir(parent)

	f, err := os.Open(root + "/test/testdata/testURL.data")
	if err != nil {
		log.Fatal("cannot get random url: ", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var urls []string
	for sc.Scan() {
		urls = append(urls, sc.Text())
	}

	return urls[rand.Intn(len(urls))]
}
