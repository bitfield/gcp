package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bitfield/gcp"
)

var g gcp.Client

func main() {
	var zone = flag.String("zone", "", "GCP zone")
	var project = flag.String("project", "", "GCP project")
	var resource = flag.String("resource", "instance", "GCP resource type (zone, instance, etc)")
	flag.Parse()
	if *project == "" {
		fmt.Println("Please specify a project with -project")
		os.Exit(1)
	}
	if err := g.Connect(); err != nil {
		log.Fatal("failed to connect to Google Cloud: %v\n", err)
	}
	switch *resource {
	case "zone":
		dumpZones(*project)
	case "instance":
		if *zone == "" {
			fmt.Println("Please specify a zone with -zone")
			os.Exit(1)
		}
		dumpInstances(*project, *zone)
	}
}

func dumpInstances(project, zone string) {
	instances, err := g.Instances(project, zone)
	if err != nil {
		log.Fatalf("failed to list instances for project %s, zone %s: %v\n", project, zone, err)
	}
	for _, i := range instances {
		json, err := i.MarshalJSON()
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

func dumpZones(project string) {
	zones, err := g.Zones(project)
	if err != nil {
		log.Fatalf("failed to list zones for project %s: %v\n", project, err)
	}
	for _, z := range zones {
		fmt.Println(z.Name)
	}
	fmt.Printf("\n\n# Retrieved %d resources\n", len(zones))
}
