package main

import (
   "io/ioutil"
   "log"
   "net/http"
   "fmt"
   "strings"
   "os"
   "encoding/json"
)

type JsonResponse struct {
	Tracks []struct {
		Album struct {
			AlbumType string `json:"album_type"`
			Artists   []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name                 string `json:"name"`
			ReleaseDate          string `json:"release_date"`
			ReleaseDatePrecision string `json:"release_date_precision"`
			TotalTracks          int    `json:"total_tracks"`
			Type                 string `json:"type"`
			URI                  string `json:"uri"`
		} `json:"album"`
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
		DiscNumber  int  `json:"disc_number"`
		DurationMs  int  `json:"duration_ms"`
		Explicit    bool `json:"explicit"`
		ExternalIds struct {
			Isrc string `json:"isrc"`
		} `json:"external_ids"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href        string `json:"href"`
		ID          string `json:"id"`
		IsLocal     bool   `json:"is_local"`
		IsPlayable  bool   `json:"is_playable"`
		Name        string `json:"name"`
		Popularity  int    `json:"popularity"`
		PreviewURL  string `json:"preview_url"`
		TrackNumber int    `json:"track_number"`
		Type        string `json:"type"`
		URI         string `json:"uri"`
	} `json:"tracks"`
}
type BearerToken struct {
	Token string `json:"token"`
}

type Token struct {
	Token string
}

func main() {
	token := getBearerToken()
	musicsIds, _ := ioutil.ReadFile("src/musics-ids.txt")
	arrayIds:=strings.Split(string(musicsIds),"\n")
	strIds:=""
	f_write_csv, _ := os.Create("src/artists-ids.txt")
	for i := 0; i < len(arrayIds); i++ {
		if (i%48 == 0 && i !=0) || i == len(arrayIds)-1 {
			resp:=makeRequest(strIds, token)
			data := JsonResponse{}
			json.Unmarshal([]byte(resp), &data)
			for x := 0; x < len(data.Tracks); x++ {
				for y := 0; y < len(data.Tracks[x].Album.Artists); y++ {
					_, err := f_write_csv.Write([]byte(data.Tracks[x].Album.Artists[y].ID+"\n"))
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			fmt.Println(len(data.Tracks))
			strIds=""
		}
		if(i%2==0){
			strIds=arrayIds[i]+strIds
		}
		strIds=arrayIds[i]+","+strIds
	}
}

func getBearerToken() string {
	url:="http://localhost:8080/token"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	 }
	body, _ := ioutil.ReadAll(resp.Body)
	token := BearerToken{}
	json.Unmarshal([]byte(string(body)), &token)
	return token.Token
}

func makeRequest(musicId string, token string) string{
	var bearer = "Bearer "+token
	url:="https://api.spotify.com/v1/tracks"
	req, err := http.NewRequest("GET", url, nil)

	q := req.URL.Query()
	q.Add("ids",musicId)
	
	req.Header.Add("Authorization", bearer)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	strJson:=string(body)
	return strJson
}