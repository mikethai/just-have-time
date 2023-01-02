package userHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type HttpClient interface {
	GetUserFollowing(msno int64, sid string) (*[]Followee, error)
	recursiveGetUserFollowing(nextPageURL string, sid string) (*[]Followee, error)
	EncryptMsno(msno int64, sid string) (*string, error)
	DecryptMsno(msno string, sid string) (*int64, error)
}

type httpClient struct {
	client *http.Client
}

func NewHttpClient(client *http.Client) *httpClient {
	return &httpClient{
		client: client,
	}
}

type MsnoEncrpytResponse struct {
	Data MsnoEncrpytResponseData `json:"data"`
}

type MsnoEncrpytResponseData struct {
	MsnoString string `json:"encryptId"`
}

type MsnoDecryptResponse struct {
	Data MsnoDecryptResponseData `json:"data"`
}

type MsnoDecryptResponseData struct {
	MsnoInt64 int64 `json:"decryptId"`
}

type Followee struct {
	Msno int64
	Type string
}

type UserFollowingResponse struct {
	Data   []UserFollowingData `json:"data"`
	Paging UserFollowingPaging `json:"paging"`
}

type UserFollowingData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type UserFollowingPaging struct {
	Next string `json:"next"`
}

func (client *httpClient) GetUserFollowing(msno int64, sid string) (*[]Followee, error) {

	var followee []Followee
	userFollowingResponse := new(UserFollowingResponse)

	msnoString, err := client.EncryptMsno(msno, sid)
	if err != nil {
		return nil, err
	}
	url := "https://api-listen-with.kkbox.com.tw/v3/users/" + *msnoString + "/following"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", sid)

	res, err := client.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("Following API return HTTP Status : " + http.StatusText(res.StatusCode))
	}

	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(reqBody, &userFollowingResponse)

	for _, user := range userFollowingResponse.Data {

		if user.Type == "user" {
			userID, _ := client.DecryptMsno(user.ID, sid)

			followee = append(followee, Followee{
				Msno: *userID,
				Type: user.Type,
			})
		}
	}

	if userFollowingResponse.Paging.Next != "" {
		// recursive to next page
		tmpFollowee, err := client.recursiveGetUserFollowing(userFollowingResponse.Paging.Next, sid)
		if err != nil {
			return nil, err
		}
		followee = append(followee, *tmpFollowee...)
	}

	return &followee, nil
}

func (client *httpClient) recursiveGetUserFollowing(nextPageURL string, sid string) (*[]Followee, error) {
	var followee []Followee
	userFollowingResponse := new(UserFollowingResponse)

	req, err := http.NewRequest("GET", nextPageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", sid)

	res, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("Get Following API return HTTP Status : " + http.StatusText(res.StatusCode))
	}

	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(reqBody, &userFollowingResponse)

	if len(userFollowingResponse.Data) > 0 {
		for _, user := range userFollowingResponse.Data {

			if user.Type == "user" {
				userID, _ := client.DecryptMsno(user.ID, sid)
				followee = append(followee, Followee{
					Msno: *userID,
					Type: user.Type,
				})
			}
		}

		tmpFollowee, _ := client.recursiveGetUserFollowing(userFollowingResponse.Paging.Next, sid)
		followee = append(followee, *tmpFollowee...)
	}

	return &followee, nil
}

func (client *httpClient) EncryptMsno(msno int64, sid string) (*string, error) {

	msnoEncrpytResponse := new(MsnoEncrpytResponse)
	url := "https://api-listen-with.kkbox.com.tw/v3/encrypt"
	method := "POST"

	payload := strings.NewReader(`{
	"id": "` + strconv.FormatInt(msno, 10) + `",
	"type": "M"
	}`)

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", sid)

	res, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("Encrypt API return HTTP Status : " + http.StatusText(res.StatusCode))
	}

	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(reqBody, &msnoEncrpytResponse)

	return &msnoEncrpytResponse.Data.MsnoString, nil
}

func (client *httpClient) DecryptMsno(msno string, sid string) (*int64, error) {

	msnoDecryptResponse := new(MsnoDecryptResponse)
	url := "https://api-listen-with.kkbox.com.tw/v3/decrypt"
	method := "POST"

	payload := strings.NewReader(`{
	"id": "` + msno + `",
	"type": "M"
	}`)

	req, _ := http.NewRequest(method, url, payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", sid)

	res, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("Decrypt API return HTTP Status : " + http.StatusText(res.StatusCode))
	}

	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(reqBody, &msnoDecryptResponse)

	return &msnoDecryptResponse.Data.MsnoInt64, nil
}
