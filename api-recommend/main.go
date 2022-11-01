package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"encoding/json"
)

type ClientResp struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Href string `json:"href"`
}

func main() {
	e := echo.New()
	route := e.Group("/recommend/api/v1")

	route.POST("/get-recommendation", func(c echo.Context) error {
		req := make(map[string]interface{})
		err := json.NewDecoder(c.Request().Body).Decode(&req)
		resp := makeRequest(req, getBearerToken())
		data := ApiResponse{}
		json.Unmarshal([]byte(resp), &data)
		
		recommendation := ClientResp{}
		var recommendations []ClientResp
		for i := 0; i < len(data.Tracks); i++ { 
			recommendation.Name = data.Tracks[i].Name
			recommendation.ID = data.Tracks[i].ID
			recommendation.Href = data.Tracks[i].Href
			recommendations = append(recommendations, recommendation)
			if err != nil {
				return c.String(http.StatusBadRequest, "error on request")
			}
		}
		return c.JSON(http.StatusOK, recommendations)
	})
	e.Logger.Fatal(e.Start(":5000"))
}