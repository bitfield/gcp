package gcp

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// Client is a wrapper for Google Cloud's various API Service types.
type Client struct {
	s   *compute.Service
	ctx context.Context
}

// New connects to Google Cloud with your application default credentials and returns a *Client ready to use
func New() (*Client, error) {
	ctx := context.Background()
	google, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("couldn't create DefaultClient: %v", err)
	}
	computeService, err := compute.New(google)
	if err != nil {
		return nil, fmt.Errorf("couldn't create compute service: %v", err)
	}
	return &Client{
		s:   computeService,
		ctx: ctx,
	}, nil
}

// Instances returns all compute instances in the specified project and zone, or nil if there was an error
func (g *Client) Instances(project, zone string) ([]*compute.Instance, error) {
	instances := []*compute.Instance{}
	instancesListCall := g.s.Instances.List(project, zone)
	if err := instancesListCall.Pages(g.ctx, func(page *compute.InstanceList) error {
		instances = append(instances, page.Items...)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("instance list call failed: %v", err)
	}
	return instances, nil
}

func (g *Client) Zones(project string) ([]string, error) {
	zones := []string{}
	zonesListCall := g.s.Zones.List(project)
	err := zonesListCall.Pages(g.ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			zones = append(zones, zone.Name)
		}
		return nil
	})
	if err != nil {
		if apiError, ok := err.(*googleapi.Error); ok {
			switch apiError.Code {
			case http.StatusForbidden:
				return nil, fmt.Errorf("Project %s is not API-enabled, skipping", project)
			case http.StatusNotFound:
				return nil, fmt.Errorf("Project %s not found", project)
			default:
				return nil, err
			}
		}
	}
	return zones, nil
}
