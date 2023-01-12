package storySongHandler

import (
	"sort"
	"time"
)

type ResponseStorySong struct {
	JID            uint      `json:"JID"`
	SongID         string    `json:"songId"`
	SongName       string    `json:"songName"`
	SongAlbumImage string    `json:"songAlbumImage"`
	Artist         string    `json:"artist"`
	SongHashTag    []string  `json:"songHashTags"`
	CreatedAt      int       `json:"createdAt"`
	CreatedTime    time.Time `json:"createdTime"`
}

type ResponseStoty struct {
	Msno           string              `json:"msno"`
	UserImage      string              `json:"userImage"`
	UserName       string              `json:"userName"`
	UserHashTags   []string            `json:"userHashTags"`
	Songs          []ResponseStorySong `json:"songs"`
	lastUpdateTime int                 `json:"lastUpdateTime"`
}

func getStorysAsSlice(storyMap map[string]ResponseStoty, msno string) []ResponseStoty {

	var storyCards []ResponseStoty
	var keys []string
	for k, v := range storyMap {
		if msno == k {
			v.lastUpdateTime = int(time.Now().Unix())
		}
		keys = append(keys, k)
		storyCards = append(storyCards, v)
	}

	sort.SliceStable(storyCards, func(i, j int) bool {
		return storyCards[i].lastUpdateTime > storyCards[j].lastUpdateTime
	})

	return storyCards
}
