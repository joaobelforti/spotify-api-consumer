package main

import (
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

func makeRequest(params map[string]interface{}, token string) string{
	var bearer = "Bearer "+token
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/recommendations", nil)
	q := req.URL.Query()
	q.Add("limit", params["limit"].(string))
	q.Add("market", params["market"].(string))
	q.Add("seed_genres", params["seed_genres"].(string))
	q.Add("min_danceability", params["min_danceability"].(string))
	q.Add("max_danceability", params["max_danceability"].(string))
	q.Add("min_energy", params["min_energy"].(string))
	q.Add("max_energy", params["max_energy"].(string))
	q.Add("min_acousticness", params["min_acousticness"].(string))
	q.Add("max_acousticness", params["max_acousticness"].(string))
	q.Add("min_instrumentalness", params["min_instrumentalness"].(string))
	q.Add("max_instrumentalness", params["max_instrumentalness"].(string))
	q.Add("min_loudness", params["min_loudness"].(string))
	q.Add("max_loudness", params["max_loudness"].(string))
	q.Add("min_speechiness", params["min_speechiness"].(string))
	q.Add("max_speechiness", params["max_speechiness"].(string))
	q.Add("min_liveness", params["min_liveness"].(string))
	q.Add("max_liveness", params["max_liveness"].(string))
	q.Add("max_valence", params["max_valence"].(string))
	q.Add("min_valence", params["min_valence"].(string))
	q.Add("max_tempo", params["max_tempo"].(string))
	q.Add("min_tempo", params["min_tempo"].(string))

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