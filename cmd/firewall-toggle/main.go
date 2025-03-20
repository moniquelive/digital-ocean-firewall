package main

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"os"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "on" && os.Args[1] != "off" {
		fmt.Println("Usage: firewall-toggle <on|off>")
		os.Exit(1)
	}
	client := godo.NewFromToken(os.Getenv("DIGITAL_OCEAN_TOKEN"))

	// get Firewalls
	list, _, err := client.Firewalls.List(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "on":
		// get Droplets
		droplets, _, err := client.Droplets.List(context.TODO(), nil)
		if err != nil {
			panic(err)
		}
		dropletsIDs := make([]int, 0, len(droplets))
		for _, d := range droplets {
			dropletsIDs = append(dropletsIDs, d.ID)
		}

		for _, r := range list {
			fmt.Printf("%v - %v - %v\n", r.DropletIDs, r.Name, r.Status)
			req := godo.FirewallRequest{
				Name:          r.Name,
				InboundRules:  r.InboundRules,
				OutboundRules: r.OutboundRules,
				DropletIDs:    dropletsIDs,
				Tags:          r.Tags,
			}
			if _, resp, err := client.Firewalls.Update(context.TODO(), r.ID, &req); err != nil {
				panic(err)
			} else {
				fmt.Println(resp)
			}
		}

	case "off":
		for _, r := range list {
			fmt.Printf("%v - %v - %v\n", r.DropletIDs, r.Name, r.Status)
			req := godo.FirewallRequest{
				Name:          r.Name,
				InboundRules:  r.InboundRules,
				OutboundRules: r.OutboundRules,
				DropletIDs:    []int{},
				Tags:          r.Tags,
			}
			if _, resp, err := client.Firewalls.Update(context.TODO(), r.ID, &req); err != nil {
				panic(err)
			} else {
				fmt.Println(resp)
			}
		}
	}
}
