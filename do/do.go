package do

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type (
	Addresses struct {
		Addresses []string `json:"addresses"`
	}
	InboundRule struct {
		Protocol string    `json:"protocol"`
		Ports    string    `json:"ports"`
		Sources  Addresses `json:"sources"`
	}
	OutboundRule struct {
		Protocol     string    `json:"protocol"`
		Ports        string    `json:"ports"`
		Destinations Addresses `json:"destinations"`
	}
	UpdateFirewall struct {
		Name          string         `json:"name"`
		InboundRules  []InboundRule  `json:"inbound_rules"`
		OutboundRules []OutboundRule `json:"outbound_rules"`
		Tags          []string       `json:"tags"`
		DropletIDs    []string       `json:"droplet_ids"`
	}
	Firewall struct {
		Id             string         `json:"id"`
		Name           string         `json:"name"`
		Status         string         `json:"status"`
		InboundRules   []InboundRule  `json:"inbound_rules"`
		OutboundRules  []OutboundRule `json:"outbound_rules"`
		CreatedAt      time.Time      `json:"created_at"`
		DropletIds     []int          `json:"droplet_ids"`
		Tags           []string       `json:"tags"`
		PendingChanges []string       `json:"pending_changes"`
	}
	Firewalls struct {
		Firewalls []Firewall `json:"firewalls"`
		Links     struct{}   `json:"links"`
		Meta      struct {
			Total int `json:"total"`
		} `json:"meta"`
	}
	DigitalOceanURL string
)

var (
	apiHeaders = http.Header{
		"Authorization": {"Bearer " + os.Getenv("DIGITAL_OCEAN_TOKEN")},
		"Content-Type":  {"application/json"},
	}
	getFirewallsURL DigitalOceanURL = "https://api.digitalocean.com/v2/firewalls"
)

func GetFirewalls() (*Firewalls, error) {
	u, err := url.Parse(string(getFirewallsURL))
	if err != nil {
		return nil, err
	}

	req := http.Request{Method: "GET", URL: u, Header: apiHeaders}
	res, err := http.DefaultClient.Do(&req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var fws Firewalls
	if err := json.NewDecoder(res.Body).Decode(&fws); err != nil {
		return nil, err
	}
	return &fws, nil
}

func PutFirewalls(firewall *Firewall) (int, error) {
	body := bytes.Buffer{}
	err := json.NewEncoder(&body).Encode(UpdateFirewall{
		Name:          firewall.Name,
		InboundRules:  firewall.InboundRules,
		OutboundRules: firewall.OutboundRules,
	})
	if err != nil {
		return -1, err
	}

	u := fmt.Sprintf("%s/%s", getFirewallsURL, firewall.Id)
	req, err := http.NewRequest("PUT", u, bytes.NewReader(body.Bytes()))
	if err != nil {
		return -1, err
	}

	req.Header = apiHeaders
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	return res.StatusCode, nil
}
