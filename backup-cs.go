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
	Acousticness     float64 `json:"acousticness"`
	AnalysisURL      string  `json:"analysis_url"`
	Danceability     float64 `json:"danceability"`
	DurationMs       int     `json:"duration_ms"`
	Energy           float64 `json:"energy"`
	ID               string  `json:"id"`
	Instrumentalness float64 `json:"instrumentalness"`
	Key              int     `json:"key"`
	Liveness         float64 `json:"liveness"`
	Loudness         float64 `json:"loudness"`
	Mode             int     `json:"mode"`
	Speechiness      float64 `json:"speechiness"`
	Tempo            float64 `json:"tempo"`
	TimeSignature    int     `json:"time_signature"`
	TrackHref        string  `json:"track_href"`
	Type             string  `json:"type"`
	URI              string  `json:"uri"`
	Valence          float64 `json:"valence"`
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
		csvLine:=fmt.Sprintf("%f", data.Danceability)+","+fmt.Sprintf("%f",data.Energy)+","+strconv.Itoa(data.Key)+","+fmt.Sprintf("%f",data.Loudness)+","+strconv.Itoa(data.Mode)+","+fmt.Sprintf("%f",data.Speechiness)+","+fmt.Sprintf("%f",data.Acousticness)+","+fmt.Sprintf("%f",data.Instrumentalness)+","+fmt.Sprintf("%f",data.Liveness)+","+fmt.Sprintf("%f",data.Valence)+","+fmt.Sprintf("%f",data.Tempo)+","+fmt.Sprintf("%f",data.Tempo)+","+string(data.Type)+","+string(data.ID)+","+string(data.URI)+","+string(data.TrackHref)+","+string(data.AnalysisURL)+","+strconv.Itoa(data.DurationMs)+","+strconv.Itoa(data.TimeSignature)+"\n"
		_, err := f.Write([]byte(csvLine))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getBearer() string {
	return "BQA4mY1ejtjhc4uqKKpo5_TbyDqZ_YXlsIrjmM7jQe4uqcfrefQvYdNgVUgLksQX6HsvGfx24Qyy16sgiiO6CjaMowmkEc_znMEne0IaltrOkLXWHJBfMHNYaH97KwXb6iB0WU9t6AxQtTTEGdAX0mJ0PruJM82WnrsgoUDOFCoP4TXBGMghjXHxDgeNVWejo7BeM7-W-wVzPiovt8GqbnSZmyaccYJI"
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
