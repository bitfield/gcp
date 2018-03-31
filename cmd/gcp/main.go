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
	project := flag.Args()[0]
	if project == "" {
		fmt.Printf("Usage: %s GCP-PROJECT\n", os.Args[0])
		os.Exit(1)
	}
	g, err := gcp.New()
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

		// Print out some metadata about the instance
		fmt.Printf(`
resource "google_compute_instance" "%s" {
  name         = "%s"
  machine_type = "%s"
  zone         = "%s"
}
`, instance.Name, instance.Name, instance.MachineType, instance.Zone)
	}

	fmt.Printf("# Retrieved %d resources\n", len(instances))
}
