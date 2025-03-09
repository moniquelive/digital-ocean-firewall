package main

import (
	"fmt"

	"github.com/moniquelive/digital-ocean-firewall/do"
	"github.com/moniquelive/digital-ocean-firewall/wtf"
)

func getIPs() (string, string, error) {
	ipv4, err := wtf.GetIP(wtf.GetIPV4URL)
	if err != nil {
		return "", "", err
	}
	ipv6, err := wtf.GetIP(wtf.GetIPV6URL)
	if err != nil {
		return "", "", err
	}
	return ipv4.IPAddress, ipv6.IPAddress, nil
}

func main() {
	// get IP's
	ipv4, ipv6, err := getIPs()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v && %v\n", ipv4, ipv6)

	// get Firewalls
	firewalls, err := do.GetFirewalls()
	if err != nil {
		panic(err)
	}
	for _, r := range firewalls.Firewalls[0].InboundRules {
		fmt.Printf("%v - %v - %v\n", r.Protocol, r.Ports, r.Sources.Addresses)
	}

	// update Firewalls
	for i, _ := range firewalls.Firewalls[0].InboundRules {
		r := &firewalls.Firewalls[0].InboundRules[i]
		if r.Ports == "22" || r.Ports == "21115" || r.Ports == "21116" || r.Ports == "21117" {
			r.Sources.Addresses = []string{ipv4, ipv6}
		}
	}

	//// print
	//bb, _ := json.MarshalIndent(firewalls, "", "  ")
	//fmt.Printf("%s\n", bb)

	// update Digital Ocean
	status, err := do.PutFirewalls(&firewalls.Firewalls[0])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Status: %d\n", status)
}
