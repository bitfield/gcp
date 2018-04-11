package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bitfield/gcp"
)

var g gcp.Client

func main() {
	var zone = flag.String("zone", "", "GCP zone")
	var project = flag.String("project", "", "GCP project")
	var resource = flag.String("resource", "", "GCP resource type (zone, instance, dnszone)")
	flag.Parse()
	if *project == "" {
		fmt.Println("Please specify a project with -project")
		os.Exit(1)
	}
	if *resource == "" {
		fmt.Println("Please specify a resource type with -resource")
		os.Exit(1)
	}
	if err := g.Connect(); err != nil {
		fmt.Printf("Failed to connect to Google Cloud: %v\n", err)
		os.Exit(1)
	}
	switch *resource {
	case "zone":
		g.ListZones(os.Stdout, *project)
	case "instance":
		checkZone(*zone)
		g.ListInstances(os.Stdout, *project, *zone)
	case "dnszone":
		g.ListDNSManagedZones(os.Stdout, *project)
	default:
		fmt.Printf("Unrecognised resource type: %s\n", *resource)
		os.Exit(1)
	}
}

func checkZone(zone string) {
	if zone == "" {
		fmt.Println("Please specify a zone with -zone")
		os.Exit(1)
	}
}
