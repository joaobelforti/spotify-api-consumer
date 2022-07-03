package main

import (
   "io/ioutil"
   "log"
   "net/http"
   "fmt"
   "strings"
   "os"
   "encoding/json"
   "strconv"
)

type MusicFeatures struct {
	Danceability     float64 `json:"danceability"`
	Energy           float64 `json:"energy"`
	Key              int     `json:"key"`
	Loudness         float64 `json:"loudness"`
	Mode             int     `json:"mode"`
	Speechiness      float64 `json:"speechiness"`
	Acousticness     float64 `json:"acousticness"`
	Instrumentalness int     `json:"instrumentalness"`
	Liveness         float64 `json:"liveness"`
	Valence          float64 `json:"valence"`
	Tempo            float64 `json:"tempo"`
	Type             string  `json:"type"`
	ID               string  `json:"id"`
	URI              string  `json:"uri"`
	TrackHref        string  `json:"track_href"`
	AnalysisURL      string  `json:"analysis_url"`
	DurationMs       int     `json:"duration_ms"`
	TimeSignature    int     `json:"time_signature"`
}

func main() {
	f, err := os.Create("musics-csv.csv")

	_, err = f.Write([]byte("danceability,energy,key,loudness,mode,speechiness,acousticness,instrumentalness,liveness,valence,tempo,type,id,uri,track_href,analysis_url,duration_ms,time_signature\n"))
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadFile("musics.txt")
    if err != nil {
        log.Fatal(err)
    }

    arrayIds:=strings.Split(string(content),"\n")

	for i := 0; i < len(arrayIds); i++ {
		resp:=makeRequest(arrayIds[i])
		data := MusicFeatures{}
		json.Unmarshal([]byte(resp), &data)
		fmt.Println(data.Danceability)
		csvLine:=fmt.Sprintf("%f", data.Danceability)+","+fmt.Sprintf("%f",data.Energy)+","+strconv.Itoa(data.Key)+","+fmt.Sprintf("%f",data.Loudness)+","+strconv.Itoa(data.Mode)+","+fmt.Sprintf("%f",data.Speechiness)+","+fmt.Sprintf("%f",data.Acousticness)+","+strconv.Itoa(data.Instrumentalness)+","+fmt.Sprintf("%f",data.Liveness)+","+fmt.Sprintf("%f",data.Valence)+","+fmt.Sprintf("%f",data.Tempo)+","+fmt.Sprintf("%f",data.Tempo)+","+string(data.Type)+","+string(data.ID)+","+string(data.URI)+","+string(data.TrackHref)+","+string(data.AnalysisURL)+","+strconv.Itoa(data.DurationMs)+","+strconv.Itoa(data.TimeSignature)+"\n"
		_, err := f.Write([]byte(csvLine))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getBearer() string {
	return "BQC7JN7G2ps6X7lsI_wLE92IiJetpPtE9GsWhl9xQ0v5H7fOBqrwa8Ecs5HVVbTGCCshNUSTbqHImCmyVkHYv1jUd28rTuCzknXg-RiAV8UnJ7akK5JumyLsX0SWIegKyI3HcZNpK2j12r_JX0XY3MH2hC7gFJDk8htcchpO-KwuL4GFPsnksg4QllEd_o1DgoFRg-dZqUBHcuggcX8firBpz2k4vebp"
}

func makeRequest(musicId string) string{
	var bearer = "Bearer "+getBearer()
	url:="https://api.spotify.com/v1/audio-features/"+musicId
	req, err := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("fields","items(track(id))")
	req.Header.Add("Authorization", bearer)
	req.URL.RawQuery = q.Encode()

	//build and send request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	//we Read the response body on the line below
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//map json response to strJson
	strJson:=string(body)
	return strJson
}
