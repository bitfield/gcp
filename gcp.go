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
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
)

// A Client holds a connected Google compute.Service and context
type Client struct {
	compute *compute.Service
	dns     *dns.Service
	ctx     context.Context
}

// Connect connects the Client to Google Cloud with your application default credentials
func (g *Client) Connect() error {
	ctx := context.Background()
	google, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return fmt.Errorf("couldn't create DefaultClient: %v", err)
	}
	g.compute, err = compute.New(google)
	if err != nil {
		return fmt.Errorf("couldn't create compute service: %v", err)
	}
	g.dns, err = dns.New(google)
	if err != nil {
		return fmt.Errorf("couldn't create DNS service: %v", err)
	}
	g.ctx = ctx
	return nil
}

// Instances returns all compute instances in the specified project and zone
func (g *Client) Instances(project, zone string) (instances []*compute.Instance, e error) {
	if err := g.compute.Instances.List(project, zone).Pages(g.ctx, func(page *compute.InstanceList) error {
		instances = append(instances, page.Items...)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to list instances for project %s, zone %s: %v", project, zone, interpretGoogleAPIError(err))
	}
	return instances, nil
}

// ListInstances gets the list of all GCP instances in the specified project and writes an HCL representation of each to the specified io.Writer.
func (g *Client) ListInstances(w io.Writer, project string, zone string) error {
	instances, err := g.Instances(project, zone)
	if err != nil {
		return err
	}
	if err = dumpInstances(w, instances); err != nil {
		return err
	}
	return nil
}

func dumpInstances(w io.Writer, instances []*compute.Instance) error {
	for _, i := range instances {
		json, err := i.MarshalJSON()
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		if err = JSON2HCL(w, json); err != nil {
			return fmt.Errorf("failed to parse JSON to HCL: %v", err)
		}
	}
	return nil
}

// Zones returns all zones in the specified project
func (g *Client) Zones(project string) (zones []*compute.Zone, e error) {
	if err := g.compute.Zones.List(project).Pages(g.ctx, func(page *compute.ZoneList) error {
		zones = append(zones, page.Items...)
		return nil
	}); err != nil {
		return nil, interpretGoogleAPIError(err)
	}
	return zones, nil
}

// ListZones gets the list of all zones in the specified project and prints their names to the specified io.Writer.
func (g *Client) ListZones(w io.Writer, project string) error {
	zones, err := g.Zones(project)
	if err != nil {
		return err
	}
	dumpZones(w, zones)
	return nil
}

func dumpZones(w io.Writer, zones []*compute.Zone) {
	for _, z := range zones {
		fmt.Fprintln(w, z.Name)
	}
}

// DNSManagedZones returns all DNS managed zones in the specified project
func (g *Client) DNSManagedZones(project string) (dnszones []*dns.ManagedZone, e error) {
	if err := g.dns.ManagedZones.List(project).Pages(g.ctx, func(page *dns.ManagedZonesListResponse) error {
		dnszones = append(dnszones, page.ManagedZones...)
		return nil
	}); err != nil {
		return nil, interpretGoogleAPIError(err)
	}
	return dnszones, nil
}

// ListDNSManagedZones gets the list of all DNS managed zones in the specified project and prints their names to the specified io.Writer.
func (g *Client) ListDNSManagedZones(w io.Writer, project string) error {
	zones, err := g.DNSManagedZones(project)
	if err != nil {
		return err
	}
	dumpDNSManagedZones(w, zones)
	return nil
}

func dumpDNSManagedZones(w io.Writer, dnszones []*dns.ManagedZone) {
	for _, z := range dnszones {
		fmt.Fprintln(w, z.DnsName)
	}
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
