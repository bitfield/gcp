package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bitfield/gcp"

	compute "google.golang.org/api/compute/v1"
)

var zone = flag.String("zone", "", "GCP zone to search")

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %s GCP-PROJECT\n", os.Args[0])
		os.Exit(1)
	}
	project := flag.Args()[0]
	g, err := gcp.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	var zones []string
	var instances []*compute.Instance
	if *zone == "" {
		zones, err = g.Zones(project)
		if err != nil {
			log.Fatalf("listing zones failed: %v\n", err)
		}
	} else {
		zones = []string{*zone}
	}
	fmt.Printf("Scanning zones: %v\n", zones)
	for _, z := range zones {
		instances, err = g.Instances(project, z)
		if err != nil {
			log.Fatalf("listing instances failed: %v\n", err)
		}
	}
	for _, instance := range instances {
		json, err := instance.MarshalJSON()
		if err != nil {
			log.Fatalf("failed to marshal JSON: %v\n", err)
		}
		hcl, err := gcp.JSON2HCL(json)
		if err != nil {
			log.Fatalf("failed to parse JSON to HCL: %v\n", err)
		}
		fmt.Println(hcl)
	}
	fmt.Printf("\n\n# Retrieved %d resources\n", len(instances))
}
