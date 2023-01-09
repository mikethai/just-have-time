package storySongHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mikethai/just-have-time/config"
)

const authHeaderField = "Authorization"
const openApiUrl = "https://api.kkbox.com/v1.1/"

type HttpClient interface {
	GetSongInfo(songID string) (*SongInfo, error)
}

type httpClient struct {
	client *http.Client
}

func NewHttpClient(client *http.Client) *httpClient {
	return &httpClient{
		client: client,
	}
}

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

func (client *httpClient) GetSongInfo(songID string) (*SongInfo, error) {
	var songInfo SongInfo

	url := openApiUrl + "tracks/" + songID + "?territory=TW"
	req, _ := http.NewRequest("GET", url, nil)
	bearerToken := config.Config("KKBOX_OPENAPI_BEARER_TOKEN")
	req.Header.Add(authHeaderField, "Bearer "+bearerToken)

	res, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(reqBody, &songInfo)

	return &songInfo, nil
}
