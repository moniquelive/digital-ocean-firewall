package main

import (
	"context"
	"fmt"
	"os"

	"github.com/moniquelive/digital-ocean-firewall/internal/wtf"

	"github.com/digitalocean/godo"
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
	fmt.Printf("%v // %v\n", ipv4, ipv6)

	// get Firewalls
	client := godo.NewFromToken(os.Getenv("DIGITAL_OCEAN_TOKEN"))
	list, _, err := client.Firewalls.List(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	// update rules
	for _, r := range list {
		fmt.Printf("%v - %v - %v\n", r.DropletIDs, r.Name, r.Status)
		req := godo.FirewallRequest{
			Name:          r.Name,
			InboundRules:  r.InboundRules,
			OutboundRules: r.OutboundRules,
			DropletIDs:    r.DropletIDs,
			Tags:          r.Tags,
		}
		rules := &req.InboundRules
		for i, _ := range *rules {
			if (*rules)[i].PortRange == "22" || (*rules)[i].PortRange == "21115" || (*rules)[i].PortRange == "21116" || (*rules)[i].PortRange == "21117" {
				(*rules)[i].Sources.Addresses = []string{ipv4, ipv6}
			}
			if _, resp, err := client.Firewalls.Update(context.TODO(), r.ID, &req); err != nil {
				panic(err)
			} else {
				fmt.Println(resp)
			}
		}
	}
}
