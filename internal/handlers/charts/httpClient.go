package chartsHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mikethai/just-have-time/config"
)

type HttpClient interface {
	GetSongCharts() (*TrackChartsResult, error)
}

type httpClient struct {
	client *http.Client
}

func NewHttpClient(client *http.Client) *httpClient {
	return &httpClient{
		client: client,
	}
}

type TrackChartsResult struct {
	Tracks TrackData `json:"tracks"`
}

type TrackData struct {
	Data []Track `json:"data"`
}

type Track struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Duration int       `json:"duration"`
	Isrc     string    `json:"isrc"`
	Url      string    `json:"url"`
	Album    AlbumInfo `json:"album"`
}

type AlbumInfo struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Url         string      `json:"url"`
	ReleaseDate string      `json:"release_date"`
	Images      []ImageInfo `json:"images"`
	Artist      ArtistInfo  `json:"artist"`
}

type ImageInfo struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Url    string `json:"url"`
}

type ArtistInfo struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	Images []ImageInfo
}

func (client *httpClient) GetSongCharts() (*TrackChartsResult, error) {
	url := "https://api.kkbox.com/v1.1/charts/LZPhK2EyYzN15dU-PT?territory=TW"

	req, _ := http.NewRequest("GET", url, nil)
	bearerToken := config.Config("KKBOX_OPENAPI_BEARER_TOKEN")
	req.Header.Add("Authorization", "Bearer "+bearerToken)

	res, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var trackChartsResult TrackChartsResult
	json.Unmarshal(reqBody, &trackChartsResult)

	return &trackChartsResult, nil
}
