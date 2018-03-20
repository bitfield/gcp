package main

import (
	"fmt"
	"log"
	"os"

	"github.com/carezone/gcp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s GCP-PROJECT\n", os.Args[0])
		os.Exit(1)
	}
	project := os.Args[1]
	g, err := gcp.New()
	if err != nil {
		log.Fatal(err)
	}
	instances, err := g.Instances(project)
	if err != nil {
		log.Fatal(err)
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
