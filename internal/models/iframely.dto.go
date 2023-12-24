package models

type Mata struct {
	URL  string `json:"url"`
	Meta meta   `json:"meta"`
	Link link   `json:"links"`
	HTML string `json:"html"`
}

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
