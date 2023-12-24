package models

import (
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

func (m *Meta) Parse(meta Mata) error {
	m.providerName(meta.URL)
	m.youtubeID(meta.URL)
	ok, err := m.htmlHasIframely(meta.HTML)
	if err != nil {
		return err
	}

	m.DataIframelyURL = ok
	m.URL = meta.URL
	m.Type = "rich"
	m.Version = "1.0"
	m.Title = meta.Meta.Title
	m.Author = meta.Meta.Author
	m.Description = meta.Meta.Description
	m.ThumbnailURL = meta.Link.Thumbnail[0].Href
	m.ThumbnailWidth = meta.Link.Thumbnail[0].Media.Width
	m.ThumbnailHeight = meta.Link.Thumbnail[0].Media.Height
	m.HTML = meta.HTML

	return nil
}

func (m *Meta) providerName(url string) {
	pat := `(?im)^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/?\n]+)`
	re := regexp.MustCompile(pat)
	match := re.FindStringSubmatch(url)

	for i := range re.SubexpNames() {
		if i != 0 {
			m.ProviderName = match[i]
		}
	}
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
