package randomutil

import (
	"bufio"
	"log"
	"math/rand"
	"os"
)

func RandURL() string {
	f, err := os.Open("testURL.data")
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
