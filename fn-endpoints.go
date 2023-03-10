package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type LightswitchResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// http://lightswitch-public-service-prod.ol.epicgames.com/lightswitch/api/service/:serviceId/status
func getLightswitchStatus(bearer string) string {
	uri := "http://lightswitch-public-service-prod.ol.epicgames.com/lightswitch/api/service/fortnite/status"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	bearerHeader := fmt.Sprint("bearer " + bearer)
	req, err := http.NewRequest("GET", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", bearerHeader)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	res := &LightswitchResponse{}

	err = json.Unmarshal(body, res)
	if err != nil {
		log.Fatalln(err)
	}

	return res.Message
}
