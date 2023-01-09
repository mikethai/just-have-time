package storySongHandler

type ResponseStorySong struct {
	SongID         string   `json:"songId"`
	SongName       string   `json:"songName"`
	SongAlbumImage string   `json:"songAlbumImage"`
	Artist         string   `json:"artist"`
	SongHashTag    []string `json:"songHashTags"`
	CreatedAt      int      `json:"createdAt"`
}

type ResponseStoty struct {
	Msno         string              `json:"msno"`
	UserImage    string              `json:"userImage"`
	UserHashTags []string            `json:"userHashTags"`
	Songs        []ResponseStorySong `json:"songs"`
}

func getStorysAsSlice(storyMap map[string]ResponseStoty, msno string) []ResponseStoty {
	// Defines the Slice length to match the Map elements count
	sm := make([]ResponseStoty, len(storyMap))

	i := 0
	for _, tx := range storyMap {
		if msno == tx.Msno {
			copy(sm[1:], sm)
			sm[0] = tx
		} else {
			sm[i] = tx
		}

		i++
	}

	return sm
}
