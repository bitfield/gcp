package gcp

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

type GcpClient struct {
	computeService *compute.Service
	ctx            context.Context
}

func New() (*GcpClient, error) {
	g := &GcpClient{}
	g.ctx = context.Background()
	c, err := google.DefaultClient(g.ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}
	g.computeService = computeService
	return g, nil
}

func (g *GcpClient) Instances(project string) ([]*compute.Instance, error) {
	instances := []*compute.Instance{}

	// Get all zones in the project
	zoneReq := g.computeService.Zones.List(project)
	if err := zoneReq.Pages(g.ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {

			// Get all instances in the zone
			//log.Printf("Searching for instances in project %s, zone %s", project, zone.Name)
			computeReq := g.computeService.Instances.List(project, zone.Name)
			if err := computeReq.Pages(g.ctx, func(page *compute.InstanceList) error {
				instances = append(instances, page.Items...)
				return nil
			}); err != nil {
				log.Fatal(err)
			}
		}
		return nil
	}); err != nil {
		if apiError, ok := err.(*googleapi.Error); ok {
			switch apiError.Code {
			case http.StatusForbidden:
				return nil, errors.New(fmt.Sprintf("Project %s is not API-enabled, skipping", project))
			case http.StatusNotFound:
				return nil, errors.New(fmt.Sprintf("Project %s not found", project))
			default:
				return nil, err
			}
		}
	}
	return instances, nil
}
