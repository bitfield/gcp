package gcp

import (
	"bytes"
	"fmt"
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

// NewClient connects to Google Cloud with your application default credentials and returns a *Client ready to use
func NewClient() (*Client, error) {
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

// Instances returns all compute instances in the specified project and zone
func (g *Client) Instances(project, zone string) (instances []*compute.Instance, e error) {
	if err := g.s.Instances.List(project, zone).Pages(g.ctx, func(page *compute.InstanceList) error {
		instances = append(instances, page.Items...)
		return nil
	}); err != nil {
		return nil, interpretGoogleAPIError(err)
	}
	return instances, nil
}

// Zones returns all zones in the specified project
func (g *Client) Zones(project string) (zones []string, e error) {
	if err := g.s.Zones.List(project).Pages(g.ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			zones = append(zones, zone.Name)
		}
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
		default:
			return fmt.Errorf("API call failed: %v", err)
		}
	}
	return err
}

func JSON2HCL(json []byte) (string, error) {
	ast, err := jsonParser.Parse(json)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON to HCL: %v\n", err)
	}
	var buf bytes.Buffer
	err = hclPrinter.Fprint(&buf, ast)
	if err != nil {
		return "", fmt.Errorf("failed to print HCL: %s\n", err)
	}
	return buf.String(), nil
}
