package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SwitchToken struct {
	AccessToken string `json:"access_token"`
}
type VerifyLink struct {
	VerifyUri  string `json:"verification_uri_complete"`
	DeviceCode string `json:"device_code"`
}
type Bearer2 struct {
	Bearer      string `json:"access_token"`
	DisplayName string `json:"displayName"`
}
type ExchangeCodeResponse struct {
	ExchangeCode string `json:"code"`
}
type BearerFromExchangeResponse struct {
	Bearer    string `json:"access_token"`
	AccountId string `json:"account_id"`
}
type DeviceAuth struct {
	AccountId string `json:"accountId"`
	DeviceId  string `json:"deviceId"`
	Secret    string `json:"secret"`
}
type ReturnDeviceAuth struct {
	AccountId   string
	DeviceId    string
	Secret      string
	DisplayName string
}
type Bearer struct {
	Bearer    string `json:"access_token"`
	AccountId string `json:"account_id"`
}

func getVerifyLink() *VerifyLink {
	bearer1 := getOauthToken()
	verifyLink := getDeviceCode(bearer1)
	return verifyLink
}
func createDeviceAuth(deviceCode string) ReturnDeviceAuth {
	bearer :=
		getBearerFromDeviceCode(deviceCode)
	exchangeCode := getExchangeCode(bearer.Bearer)
	bearer1 := getBearerFromExchange(exchangeCode)
	deviceAuth := getDeviceAuth(bearer1.Bearer, bearer1.AccountId, bearer.DisplayName)
	fmt.Println("display name: " + deviceAuth.DisplayName)
	return deviceAuth
}
func getBearerWithDeviceAuth(deviceAuth DeviceAuth) string {
	uri := "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token"
	data := url.Values{}
	data.Set("token_type", "eg1")
	data.Set("grant_type", "device_auth")
	data.Set("device_id", deviceAuth.DeviceId)
	data.Set("secret", deviceAuth.Secret)
	data.Set("account_id", deviceAuth.AccountId)
	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", "basic MzQ0NmNkNzI2OTRjNGE0NDg1ZDgxYjc3YWRiYjIxNDE6OTIwOWQ0YTVlMjVhNDU3ZmI5YjA3NDg5ZDMxM2I0MWE=")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	Bearer := &Bearer{}

	err = json.Unmarshal(body, Bearer)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode == 200 {
		return Bearer.Bearer
	}
	return "error"
}
func getDeviceCode(bearer string) *VerifyLink {
	uri := "https://account-public-service-prod03.ol.epicgames.com/account/api/oauth/deviceAuthorization"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	// data.Set("grant_type", "client_credentials")
	authHeader := fmt.Sprint("bearer " + bearer)

	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", authHeader)
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

	verifyLink := &VerifyLink{}

	err = json.Unmarshal(body, verifyLink)
	if err != nil {
		log.Fatalln(err)
	}
	return verifyLink

}
func getOauthToken() string {
	uri := "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", "basic OThmN2U0MmMyZTNhNGY4NmE3NGViNDNmYmI0MWVkMzk6MGEyNDQ5YTItMDAxYS00NTFlLWFmZWMtM2U4MTI5MDFjNGQ3")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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

	switchToken := &SwitchToken{}

	err = json.Unmarshal(body, switchToken)
	if err != nil {
		log.Fatalln(err)
	}

	return switchToken.AccessToken
}
func getBearerFromDeviceCode(deviceCode string) Bearer2 {
	println("pending get bearer from device code.")
	uri := "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "device_code")
	data.Set("device_code", deviceCode)
	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", "basic OThmN2U0MmMyZTNhNGY4NmE3NGViNDNmYmI0MWVkMzk6MGEyNDQ5YTItMDAxYS00NTFlLWFmZWMtM2U4MTI5MDFjNGQ3")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	Bearer2 := &Bearer2{}

	err = json.Unmarshal(body, Bearer2)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode == 400 {
		time.Sleep(10 * time.Second)
		val := getBearerFromDeviceCode(deviceCode)
		return val

	}
	if resp.StatusCode == 200 {
		println(Bearer2.DisplayName)
		return *Bearer2

	}
	return *Bearer2
}
func getExchangeCode(bearer string) string {
	uri := "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/exchange"
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

	ExchangeCodeResponse := &ExchangeCodeResponse{}

	err = json.Unmarshal(body, ExchangeCodeResponse)
	if err != nil {
		log.Fatalln(err)
	}

	return ExchangeCodeResponse.ExchangeCode
}
func getBearerFromExchange(exchangeCode string) BearerFromExchangeResponse {
	uri := "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "exchange_code")
	data.Set("exchange_code", exchangeCode)
	data.Set("token_type", "eg1")

	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", "basic MzQ0NmNkNzI2OTRjNGE0NDg1ZDgxYjc3YWRiYjIxNDE6OTIwOWQ0YTVlMjVhNDU3ZmI5YjA3NDg5ZDMxM2I0MWE=")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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

	Response := &BearerFromExchangeResponse{}

	err = json.Unmarshal(body, Response)
	if err != nil {
		log.Fatalln(err)
	}

	return *Response
}
func getDeviceAuth(bearer string, accountId string, displayName string) ReturnDeviceAuth {
	fmt.Println("getting device auth")
	uri := fmt.Sprint("https://account-public-service-prod.ol.epicgames.com/account/api/public/account/" + accountId + "/deviceAuth")
	data := url.Values{}

	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	bearerHeader2 := fmt.Sprint("bearer " + bearer)
	req.Header.Set("Authorization", bearerHeader2)
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

	deviceAuth := &DeviceAuth{}

	err = json.Unmarshal(body, deviceAuth)
	if err != nil {
		log.Fatalln(err)
	}
	println(deviceAuth.AccountId)
	newDeviceAuth := &ReturnDeviceAuth{Secret: deviceAuth.Secret, DeviceId: deviceAuth.DeviceId, AccountId: deviceAuth.AccountId, DisplayName: displayName}
	return *newDeviceAuth
}
