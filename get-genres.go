package main

import (
   "io/ioutil"
   "log"
   "net/http"
   "strings"
   "os"
   "encoding/json"
   "strconv"
)

type JsonResponse struct {
	Artists []struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Followers struct {
			Href  interface{} `json:"href"`
			Total int         `json:"total"`
		} `json:"followers"`
		Genres []string `json:"genres"`
		Href   string   `json:"href"`
		ID     string   `json:"id"`
		Images []struct {
			Height int    `json:"height"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name       string `json:"name"`
		Popularity int    `json:"popularity"`
		Type       string `json:"type"`
		URI        string `json:"uri"`
	} `json:"artists"`
}

type BearerToken struct {
	Token string `json:"token"`
}

type Genres struct {
	genres []Genre
}
type Genre struct {
	name string
	frequency int
}

func main() {
	token := getBearerToken()
	musicsIds, _ := ioutil.ReadFile("src/artists-ids.txt")
	arrayIds:=strings.Split(string(musicsIds),"\n")
	strIds:=""
	f_write_csv, _ := os.Create("src/genres.txt")
	for i := 0; i < len(arrayIds); i++ {
		if (i%48 == 0 && i !=0) || i == len(arrayIds)-1 {
			resp:=makeRequest(strIds, token)
			strIds=""
			data := JsonResponse{}
			json.Unmarshal([]byte(resp), &data)
				for x := 0; x < len(data.Artists); x++ {
					for y:=0; y < len(data.Artists[x].Genres); y++ {
					_, err := f_write_csv.Write([]byte(data.Artists[x].Genres[y]+"\n"))
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
		if(i%2==0){
			strIds=arrayIds[i]+strIds
		}
		strIds=arrayIds[i]+","+strIds
	}
	orderGenres()
}

func orderGenres() {
	genres, _ := ioutil.ReadFile("src/genres.txt")
	arrayGenres:=strings.Split(string(genres),"\n")
    printUniqueValue(arrayGenres)
}

func printUniqueValue( arr []string){
    freq := make(map[string]int)
    for _ , num :=  range arr {
        freq[num] = freq[num]+1
    }
	genres:=Genres{}
	for value, key :=  range freq{
		genre:=Genre{}
		genre.name=value
		genre.frequency=key
		genres.genres = append(genres.genres,genre)
	}

	for x:=0; x<len(genres.genres)-1; x++{
		for y:=0; y<len(genres.genres)-1; y++{
			genreAux:=Genre{}
			if genres.genres[y].frequency > genres.genres[y+1].frequency {
				genreAux=genres.genres[y]
				genres.genres[y]=genres.genres[y+1]
				genres.genres[y+1]=genreAux
			}
		}
	}
	f_write_csv, _ := os.Create("src/order-genres.csv")
	f_write_csv.Write([]byte("genero"+","+"quantidade"+"\n"))
	for x:=0; x<len(genres.genres)-1; x++{
		f_write_csv.Write([]byte(genres.genres[x].name+","+strconv.Itoa(genres.genres[x].frequency)+"\n"))
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
	url:="https://api.spotify.com/v1/artists"
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