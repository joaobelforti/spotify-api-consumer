package main

import (
   "io/ioutil"
   "log"
   "net/http"
   "fmt"
   "strings"
   "os"
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
	strIds:=""
	f1, err1 := os.Create("test.txt")
	if err1 != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(arrayIds)+1; i++ {

		if i%100 == 0 || i==len(arrayIds)-1 && i !=0 {
			resp:=makeRequest(strIds)
			fmt.Println(resp)
			_, err1 := f1.Write([]byte(resp))
			if err1 != nil {
				log.Fatal(err)
			}
			if (i==len(arrayIds)-1){
				break;
			}
			strIds=""
		}
		strIds=arrayIds[i]+","+strIds
	}
}
func buildIdsStringForRequest(){

}

func getBearer() string {
	return "BQD0Mm4bRqYajV8hKTA4NvgXI_8ai2EgIiOnfN5Ol_3pX156dOl4uvykpOJHDXWoW6nrjFHQsE38Wzr0O0HFMnpZ_2OzDt7mqZJ9MInbGczJ67MFZKmKT4mPkWLTKhAwB0dTzIRF1FfVzOEBa7M8ak5l0MThjaZjCu5qu3-croHQAvJBDEpiXj7-bN6YessaJtAXamIJkNHJJeDUb7QVZPjyZVOCrQV5"
}

func makeRequest(musicId string) string{
	var bearer = "Bearer "+getBearer()
	url:="https://api.spotify.com/v1/audio-features"
	req, err := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("ids",musicId)
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
