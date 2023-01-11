package model

type SongInfo struct {
	Id       string    `json:"id" firestore:"id,omitempty"`
	Name     string    `json:"name" firestore:"name,omitempty"`
	Duration int       `json:"duration" firestore:"duration,omitempty"`
	Isrc     string    `json:"isrc" firestore:"isrc,omitempty"`
	Url      string    `json:"url" firestore:"url,omitempty"`
	Album    AlbumInfo `json:"album" firestore:"album,omitempty"`
}

type AlbumInfo struct {
	Id          string      `json:"id" firestore:"id,omitempty"`
	Name        string      `json:"name" firestore:"name,omitempty"`
	Url         string      `json:"url" firestore:"url,omitempty"`
	ReleaseDate string      `json:"release_date" firestore:"release_date,omitempty"`
	Images      []ImageInfo `json:"images" firestore:"images,omitempty"`
	Artist      ArtistInfo  `json:"artist" firestore:"artist,omitempty"`
}

type ImageInfo struct {
	Height int    `json:"height" firestore:"height,omitempty"`
	Width  int    `json:"width" firestore:"width,omitempty"`
	Url    string `json:"url" firestore:"url,omitempty"`
}

type ArtistInfo struct {
	Id     string      `json:"id" firestore:"id,omitempty"`
	Name   string      `json:"name" firestore:"name,omitempty"`
	Url    string      `json:"url" firestore:"url,omitempty"`
	Images []ImageInfo `json:"images" firestore:"images,omitempty"`
}
