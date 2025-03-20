package wtf

import (
	"encoding/json"
	"net/http"
)

type (
	IPs struct {
		IPAddress   string `json:"YourFuckingIPAddress"`
		Location    string `json:"YourFuckingLocation"`
		Hostname    string `json:"YourFuckingHostname"`
		ISP         string `json:"YourFuckingISP"`
		TorExit     bool   `json:"YourFuckingTorExit"`
		City        string `json:"YourFuckingCity"`
		Country     string `json:"YourFuckingCountry"`
		CountryCode string `json:"YourFuckingCountryCode"`
	}
	getIPURL string
)

var (
	GetIPV4URL getIPURL = "https://ipv4.myip.wtf/json"
	GetIPV6URL getIPURL = "https://myip.wtf/json"
)

func GetIP(u getIPURL) (*IPs, error) {
	res, err := http.DefaultClient.Get(string(u))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var ips IPs
	if err := json.NewDecoder(res.Body).Decode(&ips); err != nil {
		return nil, err
	}
	return &ips, nil
}
