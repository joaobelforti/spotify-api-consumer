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

func main() {
	f, err := os.Create("musics.txt")
	if err != nil {
		log.Fatal(err)
	}
	total:=getTotal()/100
	total=total+1
	for i := 0; i < total; i++ {
		resp:=makeRequest(i*100)
		arrayTracks:=processResponse(resp)
		for x := 0; x < len(arrayTracks); x++ {
			data := TrackId{}
			json.Unmarshal([]byte(arrayTracks[x]), &data)
			if data.Track.ID != "" {
				_, err := f.Write([]byte(data.Track.ID+"\n"))
				if err != nil {
					log.Fatal(err)
				}
			}
			
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("MUSICS IDS DONE.")
}

func processResponse(resp  string) []string{
	arrayTracks:=strings.Split(resp,"[")
	remove := regexp.MustCompile(`]`)
	arrayIds:=remove.Split(arrayTracks[1],-1)
	arrayTracks=strings.Split(arrayIds[0],",")
	return arrayTracks
}
func getBearer() string {
	return "BQDA6BQJ4WGwVNloL6RuR_CR6w4tFW-D-7DE33uo7GED_eKxIvjev7YPhTI735PYckIMkNBjQLjpwJ3pUZJwlj0Kv5WnsqERzmHDJD9NP5l0RfL4DoyOQh57QVwE9ZowhQPpR_Hq07wcOZU1WZ77t-X2ZVVnhOX8iqwIjOQVad6nYzts76j7Md45zAQJ7e9cc7EbZsqls8h5YaIPDoths6SIzN3rDjhH"
}
//https://open.spotify.com/playlist/7y6tlnIjgyIQPjmYhauphK?si=573ea74fbe2d40af
func makeRequest(offset int) string{
	var bearer = "Bearer "+getBearer()
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/37i9dQZF1DWWQRwui0ExPn/tracks", nil)
	q := req.URL.Query()
    q.Add("offset", strconv.Itoa(offset))
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

func getTotal() int {
	var bearer = "Bearer "+getBearer()
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/37i9dQZF1DWWQRwui0ExPn/tracks", nil)
	q := req.URL.Query()
	q.Add("fields","total")
	req.Header.Add("Authorization", bearer)
	req.URL.RawQuery = q.Encode()
	
	//build and send request
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