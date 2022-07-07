package main

import (
   "io/ioutil"
   "log"
   "net/http"
   "fmt"
   "strings"
   "encoding/json"
   "strconv"
   "time"
   "os"
   "sync"
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

type BearerToken struct {
	Token string `json:"token"`
}

type Token struct {
	Token string
}

func main() {
	start := time.Now()
	token := getBearerToken()
	f_write_csv, _ := os.Create("src/musics-csv.csv")
	musicsIds, _ := ioutil.ReadFile("src/musics-ids.txt")
	f_write_csv.Write([]byte("danceability,energy,key,loudness,mode,speechiness,acousticness,instrumentalness,liveness,valence,tempo,type,id,uri,track_href,analysis_url,duration_ms,time_signature\n"))
    arrayIds:=strings.Split(string(musicsIds),"\n")
	strIds:=""
	for i := 0; i < len(arrayIds)+1; i++ {
		if (i%100 == 0 || i==len(arrayIds)-1) && i !=0 {
			f1, _ := os.Create("src/tmp.txt")
			resp:=makeRequest(strIds, token)
			f1.Write([]byte(resp))
			writeCsv(f_write_csv)
			if (i==len(arrayIds)-1){
				break;
			}
			strIds=""
		}
		strIds=arrayIds[i]+","+strIds
	}
	elapsed := time.Since(start)
    log.Printf("time took = %s", elapsed)
	fmt.Println("CSV DONE.")
}

func processResponse(audio_features_read string) []string{
	audio_features_array_aux1 := strings.Split(audio_features_read,"[")
	audio_features_array_aux2 := strings.Split(audio_features_array_aux1[1],"]")
	res1 := strings.ReplaceAll(audio_features_array_aux2[0], "},{", "} , {")
	jsonArray := strings.Split(res1,", ")
	return jsonArray
}

func writeCsv(f_write_csv *os.File) {
	audio_features_read, _ := ioutil.ReadFile("src/tmp.txt")
	jsonArray := processResponse(string(audio_features_read))

	var wg sync.WaitGroup
	wg.Add(len(jsonArray))

	for i := 0; i < len(jsonArray); i++ {
		go func(i int) {
			data := MusicFeatures{}
			json.Unmarshal([]byte(jsonArray[i]), &data)
			csvLine:=fmt.Sprintf("%f", data.Danceability)+","+fmt.Sprintf("%f",data.Energy)+","+strconv.Itoa(data.Key)+","+fmt.Sprintf("%f",data.Loudness)+","+strconv.Itoa(data.Mode)+","+fmt.Sprintf("%f",data.Speechiness)+","+fmt.Sprintf("%f",data.Acousticness)+","+strconv.Itoa(data.Instrumentalness)+","+fmt.Sprintf("%f",data.Liveness)+","+fmt.Sprintf("%f",data.Valence)+","+fmt.Sprintf("%f",data.Tempo)+","+string(data.Type)+","+string(data.ID)+","+string(data.URI)+","+string(data.TrackHref)+","+string(data.AnalysisURL)+","+strconv.Itoa(data.DurationMs)+","+strconv.Itoa(data.TimeSignature)+"\n"
			_, err := f_write_csv.Write([]byte(csvLine))
			if err != nil {
				log.Fatal(err)
			}
			defer wg.Done()
		}(i)
	}
	wg.Wait()
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
	url:="https://api.spotify.com/v1/audio-features"
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