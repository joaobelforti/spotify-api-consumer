package main

import (
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

func makeRequest(params map[string]interface{}, token string) string{
	var bearer = "Bearer "+token
	fmt.Println(bearer)
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/recommendations", nil)
	q := req.URL.Query()
	fmt.Println(params["seed_tracks"])
	q.Add("limit", params["limit"].(string))
	q.Add("seed_tracks", params["seed_tracks"].(string))
	
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
	url:="http://token-container:3000/token"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	 }
	body, _ := ioutil.ReadAll(resp.Body)
	token := BearerToken{}
	json.Unmarshal([]byte(string(body)), &token)
	return token.Token
}