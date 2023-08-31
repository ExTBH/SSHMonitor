package geoip

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type GeoIP struct {
	IP          string
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	CityName    string `json:"cityName"`
}

// attemps to get the Ip info and returns an http.StatusCode
func (ip *GeoIP) Get() int {

	url := fmt.Sprintf("https://freeipapi.com/api/json/%s", ip.IP)

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode
	}
	defer resp.Body.Close()

	buff, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(buff, ip)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
