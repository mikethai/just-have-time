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

const authHeaderField = "Authorization"
const openApiUrl = "https://api.kkbox.com/v1.1/"
const listenWithApiUrl = "https://api-listen-with.kkbox.com.tw/v3/"

type HttpClient interface {
	GetUserFollowing(msno string, sid string) (*[]Followee, error)
	GetUserProfile(msno string, sid string) (*UserProfileData, error)
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
	Msno string
	Name string
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

type UserProfileResponse struct {
	Data UserProfileData `json:"data"`
}

type UserProfileData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (client *httpClient) GetUserFollowing(msno string, sid string) (*[]Followee, error) {

	var followee []Followee
	userFollowingResponse := new(UserFollowingResponse)

	url := listenWithApiUrl + "users/" + msno + "/following"

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
			followee = append(followee, Followee{
				Msno: user.ID,
				Name: user.Name,
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

func (client *httpClient) GetUserProfile(msno string, sid string) (*UserProfileData, error) {

	userProfileResponse := new(UserProfileResponse)

	url := listenWithApiUrl + "users/" + msno

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(authHeaderField, sid)

	res, err := client.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("User API return HTTP Status : " + http.StatusText(res.StatusCode))
	}

	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(reqBody, &userProfileResponse)

	return &userProfileResponse.Data, nil
}

func (client *httpClient) recursiveGetUserFollowing(nextPageURL string, sid string) (*[]Followee, error) {
	var followee []Followee
	userFollowingResponse := new(UserFollowingResponse)

	req, err := http.NewRequest("GET", nextPageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(authHeaderField, sid)

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
				followee = append(followee, Followee{
					Msno: user.ID,
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
	url := listenWithApiUrl + "encrypt"
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
	req.Header.Add(authHeaderField, sid)

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
	url := listenWithApiUrl + "decrypt"
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
