package model

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Meta struct {
	URL             string        `json:"url"`
	Type            string        `json:"type"`
	Version         string        `json:"version"`
	Title           string        `json:"title"`
	Author          string        `json:"author"`
	ProviderName    string        `json:"provider_name"`
	Description     string        `json:"description"`
	YouTubeID       string        `json:"youtube_video_id,omitempty"`
	ThumbnailURL    string        `json:"thumbnail_url"`
	ThumbnailWidth  int           `json:"thumbnail_width"`
	ThumbnailHeight int           `json:"thumbnail_height"`
	HTML            string        `json:"html"`
	CacheAge        time.Duration `json:"cache_age"`
	DataIframelyURL bool          `json:"data_iframely_url"`
}

func (m *Meta) ParseJSON(data []byte) error {
	type media struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}

	type thumbnail struct {
		Href  string `json:"href"`
		Media media  `json:"media"`
	}
	type link struct {
		Thumbnail []thumbnail `json:"thumbnail"`
	}
	type meta struct {
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
	}

	type iframely struct {
		URL  string `json:"url"`
		Meta meta   `json:"meta"`
		Link link   `json:"links"`
		HTML string `json:"html"`
	}

	mt := iframely{}
	err := json.Unmarshal(data, &mt)
	if err != nil {
		return err
	}

	tmp := &Meta{
		URL:             mt.URL,
		Type:            "rich",
		Version:         "1.0",
		Title:           mt.Meta.Title,
		Author:          mt.Meta.Author,
		Description:     mt.Meta.Description,
		ThumbnailURL:    mt.Link.Thumbnail[0].Href,
		ThumbnailWidth:  mt.Link.Thumbnail[0].Media.Width,
		ThumbnailHeight: mt.Link.Thumbnail[0].Media.Height,
		HTML:            mt.HTML,
	}

	*m = *tmp
	m.providerName()
	m.youtubeID()

	ok, err := m.htmlHasIframely()
	if err != nil {
		return err
	}
	m.DataIframelyURL = ok
	return nil
}

func (m *Meta) providerName() {
	pat := `(?im)^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/?\n]+)`
	re := regexp.MustCompile(pat)
	match := re.FindStringSubmatch(m.URL)

	for i := range re.SubexpNames() {
		if i != 0 {
			m.ProviderName = match[i]
		}
	}
}

func (m *Meta) youtubeID() {
	pat := `^.*(?:(?:youtu\.be\/|v\/|vi\/|u\/\w\/|embed\/|shorts\/)|(?:(?:watch)?\?v(?:i)?=|\&v(?:i)?=))([^#\&\?]*).*`
	re := regexp.MustCompile(pat)
	match := re.FindStringSubmatch(m.URL)
	if len(match) != 2 {
		return
	} else {
		m.YouTubeID = match[1]
	}
}

func (m *Meta) htmlHasIframely() (bool, error) {
	root, err := html.Parse(strings.NewReader(m.HTML))
	if err != nil {
		return false, err
	}

	ok := elementHasID("data-iframely-url", root)
	if !ok {
		return false, nil
	}
	return true, nil
}

func elementHasID(id string, n *html.Node) bool {
	for _, a := range n.Attr {
		if a.Key == id {
			return true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if ok := elementHasID(id, c); ok {
			return true
		}
	}
	return false
}
