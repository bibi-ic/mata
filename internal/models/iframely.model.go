package models

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type Meta struct {
	URL             string `json:"url"`
	Type            string `json:"type"`
	Version         string `json:"version"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	ProviderName    string `json:"provider_name"`
	Description     string `json:"description"`
	YouTubeID       string `json:"youtube_video_id,omitempty"`
	ThumbnailURL    string `json:"thumbnail_url"`
	ThumbnailWidth  int    `json:"thumbnail_width"`
	ThumbnailHeight int    `json:"thumbnail_height"`
	HTML            string `json:"html"`
	CacheAge        int64  `json:"cache_age"`
	DataIframelyURL bool   `json:"data_iframely_url"`
}

func (m *Meta) Parse() error {
	m.youtubeID(m.URL)
	ok, err := m.htmlHasIframely(m.HTML)
	if err != nil {
		return err
	}

	m.DataIframelyURL = ok
	return nil
}

func (m *Meta) youtubeID(string) {
	pat := `^.*(?:(?:youtu\.be\/|v\/|vi\/|u\/\w\/|embed\/|shorts\/)|(?:(?:watch)?\?v(?:i)?=|\&v(?:i)?=))([^#\&\?]*).*`
	re := regexp.MustCompile(pat)
	match := re.FindStringSubmatch(m.URL)
	if len(match) != 2 {
		return
	} else {
		m.YouTubeID = match[1]
	}
}

func (m *Meta) htmlHasIframely(h string) (bool, error) {
	root, err := html.Parse(strings.NewReader(h))
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
