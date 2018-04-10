package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bitfield/gcp"
	compute "google.golang.org/api/compute/v1"
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
		log.Fatalf("failed to connect to Google Cloud: %v\n", err)
	}
	switch *resource {
	case "zone":
		dumpZones(*project)
	case "instance":
		if *zone == "" {
			fmt.Println("Please specify a zone with -zone")
			os.Exit(1)
		}
		instances, err := g.Instances(*project, *zone)
		if err != nil {
			log.Fatal(err)
		}
		if err = dumpInstances(os.Stdout, instances); err != nil {
			log.Fatal(err)
		}
	}
}

func dumpInstances(w io.Writer, instances []*compute.Instance) error {
	for _, i := range instances {
		json, err := i.MarshalJSON()
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		if err = gcp.JSON2HCL(w, json); err != nil {
			return fmt.Errorf("failed to parse JSON to HCL: %v", err)
		}
	}
	return nil
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
