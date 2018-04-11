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
		log.Fatalf("failed to connect to Google Cloud: %v\n", err)
	}
	switch *resource {
	case "zone":
		g.ListZones(os.Stdout, *project)
	case "instance":
		if *zone == "" {
			fmt.Println("Please specify a zone with -zone")
			os.Exit(1)
		}
		g.ListInstances(os.Stdout, *project, *zone)
	}
}
