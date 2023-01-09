package chartsHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mikethai/just-have-time/config"
	"github.com/mikethai/just-have-time/internal/firestoreClient"
)

const authHeaderField = "Authorization"
const openApiUrl = "https://api.kkbox.com/v1.1/"
const chartPlaylistId = "LZPhK2EyYzN15dU-PT"

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
	var trackChartsResult TrackChartsResult

	url := openApiUrl + "charts/" + chartPlaylistId + "?territory=TW"

	firestoreCache := firestoreClient.NewFirestoreClient()
	defer firestoreCache.CloseConnection()
	dsnap, err := firestoreCache.Get("chart", chartPlaylistId)
	if err != nil {
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

		json.Unmarshal(reqBody, &trackChartsResult)

		firestoreCache.Set("chart", chartPlaylistId, &trackChartsResult)
		return &trackChartsResult, nil
	}

	dsnap.DataTo(&trackChartsResult)

	return &trackChartsResult, nil
}
