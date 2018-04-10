package gcp

import (
	"fmt"
	"io"
	"net/http"

	hclPrinter "github.com/hashicorp/hcl/hcl/printer"
	jsonParser "github.com/hashicorp/hcl/json/parser"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// A Client holds a connected Google compute.Service and context
type Client struct {
	s   *compute.Service
	ctx context.Context
}

// Connect connects the Client to Google Cloud with your application default credentials
func (g *Client) Connect() error {
	ctx := context.Background()
	google, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return fmt.Errorf("couldn't create DefaultClient: %v", err)
	}
	computeService, err := compute.New(google)
	if err != nil {
		return fmt.Errorf("couldn't create compute service: %v", err)
	}
	g.s = computeService
	g.ctx = ctx
	return nil
}

// Instances returns all compute instances in the specified project and zone
func (g *Client) Instances(project, zone string) (instances []*compute.Instance, e error) {
	if err := g.s.Instances.List(project, zone).Pages(g.ctx, func(page *compute.InstanceList) error {
		instances = append(instances, page.Items...)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to list instances for project %s, zone %s: %v", project, zone, interpretGoogleAPIError(err))
	}
	return instances, nil
}

// Zones returns all zones in the specified project
func (g *Client) Zones(project string) (zones []*compute.Zone, e error) {
	if err := g.s.Zones.List(project).Pages(g.ctx, func(page *compute.ZoneList) error {
		zones = append(zones, page.Items...)
		return nil
	}); err != nil {
		return nil, interpretGoogleAPIError(err)
	}
	return zones, nil
}

func interpretGoogleAPIError(err error) error {
	if apiError, ok := err.(*googleapi.Error); ok {
		switch apiError.Code {
		case http.StatusForbidden:
			return fmt.Errorf("project is not API-enabled")
		case http.StatusNotFound:
			return fmt.Errorf("project not found")
		case http.StatusBadRequest:
			return fmt.Errorf("zone not found")
		default:
			return fmt.Errorf("API call failed: %v", err)
		}
	}
	return err
}

// JSON2HCL takes a JSON representation of a GCP resource and writes the equivalent HCL (Terraform) representation to the supplied io.Writer
func JSON2HCL(w io.Writer, json []byte) error {
	ast, err := jsonParser.Parse(json)
	if err != nil {
		return fmt.Errorf("failed to parse JSON to HCL: %v", err)
	}
	err = hclPrinter.Fprint(w, ast)
	if err != nil {
		return fmt.Errorf("failed to print HCL: %s", err)
	}
	return nil
}
