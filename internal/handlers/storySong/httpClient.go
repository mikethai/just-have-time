package storySongHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mikethai/just-have-time/config"
	"github.com/mikethai/just-have-time/internal/model"
)

const authHeaderField = "Authorization"
const openApiUrl = "https://api.kkbox.com/v1.1/"

type HttpClient interface {
	GetSongInfo(songID string) (*model.SongInfo, error)
}

type httpClient struct {
	client *http.Client
}

func NewHttpClient(client *http.Client) *httpClient {
	return &httpClient{
		client: client,
	}
}

func (client *httpClient) GetSongInfo(songID string) (*model.SongInfo, error) {
	var songInfo model.SongInfo

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
