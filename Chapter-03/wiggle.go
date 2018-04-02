package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	user = "<User>"
	pass = "<Pass>"
	url  = "https://api.wigle.net/api/v2/network/search"
)

type wiggleNetwork struct {
	trilat  float64
	trilong float64
	ssid    string
	city    string
	country string
}

type wiggleResponse struct {
	success      bool
	totalResults int
	resultCount  int
	results      []wiggleNetwork
}

// GeoCode contains geo location for a given mac address
type GeoCode struct {
	ssid     string
	lat, lng float64
}

// WiggleGet fetches location information from wiggle.net for given mac address
func WiggleGet(mac string) (GeoCode, error) {
	// Make request
	client := &http.Client{}
	request, err := http.NewRequest("GET", url+"?netid="+mac, nil)
	request.SetBasicAuth(user, pass)
	response, err := client.Do(request)
	if err != nil {
		return GeoCode{}, err
	}
	defer response.Body.Close()

	// Decode response
	var wiggleresponse wiggleResponse
	if err := json.NewDecoder(response.Body).Decode(wiggleresponse); err != nil {
		return GeoCode{}, err
	}
	if !wiggleresponse.success || wiggleresponse.resultCount == 0 {
		return GeoCode{}, errors.New("could not fetch details")
	}
	result := wiggleresponse.results[0]

	// Send output
	return GeoCode{
		lat:  result.trilat,
		lng:  result.trilong,
		ssid: result.ssid,
	}, nil
}
