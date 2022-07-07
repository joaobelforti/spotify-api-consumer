package main

import (
   "io/ioutil"
   "log"
   "net/http"
   "strconv"
   "strings"
   "encoding/json"
   "regexp"
   "os"
   "fmt"
   "sync"
   "time"
)

type TrackId struct {
	Track struct {
		ID string `json:"id"`
	} `json:"track"`
}

type Total struct {
	Total int `json:"total"`
}

type BearerToken struct {
	Token string `json:"token"`
}

func main() {
	token:=getBearerToken()
	start := time.Now()
	playlists:=getPlaylists()
	f, err := os.Create("src/musics-ids.txt")
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for m := 0; m < len(playlists); m++ {
		total:=getTotalMusicsPlaylist(playlists[m], token)/100
		total=total+1

		for i := 0; i < total; i++ {
				resp:=makeRequest(i*100, playlists[m], token)
				arrayTracks:=processResponse(resp)

				wg.Add(len(arrayTracks))
				for x := 0; x < len(arrayTracks); x++ {
					go func(x int) {
						data := TrackId{}
						json.Unmarshal([]byte(arrayTracks[x]), &data)
						if data.Track.ID != "" {
							_, err := f.Write([]byte(data.Track.ID+"\n"))
							if err != nil {
								log.Fatal(err)
							}
						}
						defer wg.Done()
					}(x)
				}
			}
		}
	wg.Wait()
	elapsed := time.Since(start)
    log.Printf("time took = %s", elapsed)
	fmt.Println("MUSICS IDS DONE.")
}


func processResponse(resp  string) []string{
	arrayTracks:=strings.Split(resp,"[")
	remove := regexp.MustCompile(`]`)
	arrayIds:=remove.Split(arrayTracks[1],-1)
	arrayTracks=strings.Split(arrayIds[0],",")
	return arrayTracks
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

func makeRequest(offset int, playlistId string, token string) string{
	var bearer = "Bearer "+token
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+playlistId+"/tracks", nil)
	q := req.URL.Query()
    q.Add("offset", strconv.Itoa(offset))
	q.Add("fields","items(track(id))")
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

func getPlaylists() []string{
	musicsIds, _ := ioutil.ReadFile("playlists.txt")
	arrayMusicsIds := strings.Split(string(musicsIds),"\n")
	arraysPlaylists := []string{}
	for i := 0; i < len(arrayMusicsIds); i++ {
		reg := regexp.MustCompile(`https://open.spotify.com/playlist/`)
		textAux1 := reg.Split(arrayMusicsIds[i],-1)
		textAux2 := strings.Split(textAux1[1],"?")
		arraysPlaylists = append(arraysPlaylists, textAux2[0])
	}
	return arraysPlaylists
}

func getTotalMusicsPlaylist(playlistId  string, token string) int {
	var bearer = "Bearer "+token

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+playlistId+"/tracks", nil)
	q := req.URL.Query()
	q.Add("fields","total")

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

	data := Total{}
    json.Unmarshal([]byte(body), &data)
	return int(data.Total)
}