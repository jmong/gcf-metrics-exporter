package metricsexporter
/**
 * Client to retrieve information about GCP compute resources.
 * 
 * @see REST call definitions - https://cloud.google.com/compute/docs/reference/rest/v1/
 **/

import (
    "fmt"
    "context"
    "encoding/json"

    compute  "google.golang.org/api/compute/v1"
)

type ComputeBuilder interface {
    Context(context.Context)  ComputeBuilder
    Project(string)           ComputeBuilder
    Region(string)            ComputeBuilder
    Zone(string)              ComputeBuilder
    EnableEmitter()           ComputeBuilder

    Build()                   (Compute, error)
}

type computeBuild struct {
    context        context.Context
    project        string
    region         string
    zone           string
    enableemitter  bool
}

/* Compute object.
 */
type Compute struct {
    GcpMetadata
    context        context.Context
    client         *compute.Service
    Project        string
    Region         string
    Zone           string
    EnableEmitter  bool
    emitter        Emitters
}

/* NewComputeBuilder creates a builder object by adding components/features 
 * that will create a Compute object.
 */
func NewComputeBuilder() ComputeBuilder {
    return &computeBuild{}
}

/* Context is the Google background context of the request.
 */
func (b *computeBuild) Context(ctx context.Context) ComputeBuilder {
	b.context = ctx
	return b
}

/* Project is the GCP project id.
 */
func (b *computeBuild) Project(project string) ComputeBuilder {
	b.project = project
	return b
}

/* Region is the GCP region the Computeing resource resides in.
 */
func (b *computeBuild) Region(region string) ComputeBuilder {
	b.region = region
	return b
}

/* Zone is the GCP region the Computeing resource resides in.
 */
func (b *computeBuild) Zone(zone string) ComputeBuilder {
	b.zone = zone
	return b
}

/*
 */
func (b *computeBuild) EnableEmitter() ComputeBuilder {
	b.enableemitter = true
	return b
}

/* Build creates a Compute object that retrieves information
 * about the GCP Computeing infrastructure.
 */
func (b *computeBuild) Build() (Compute, error) {
    client, err := compute.NewService(b.context)
	if err != nil {
		return Compute{}, err
	}

    var pusher *PrometheusPush = nil
    if b.enableemitter == true {
        pusher = NewPrometheusPush(PROM_PUSHGW_URL, PROM_PUSHGW_JOB)
    }

    return Compute{
        context:        b.context,
        client:         client,
        Project:        b.project,
        Region:         b.region,
        Zone:           b.zone,
        EnableEmitter:  b.enableemitter,
        emitter:        pusher,
    }, nil
}

/* @TODO
 * Do acts on your request to retrieve and return a response to you.
 */
func (n *Compute) Do(qry Query) (string, error) {
    if qry.Resource == "compute" && qry.Action == "get" && qry.Target == "regions.list" {
        return n.getRegionsList()
    } else if qry.Resource == "compute" && qry.Action == "get" && qry.Target == "instances.list" {
        return n.getInstancesList()
    }
    return "[Debug] It will call some Compute operations to return json response", nil
}

/* @see https://cloud.google.com/compute/docs/reference/rest/v1/subComputes/list
 */
func (n *Compute) getRegionsList() (string, error) {
    var res string
    var svc = n.client
    
    list, err := svc.Regions.List(n.Project).Do()
	if err != nil {
		return fmt.Sprintf("failed to list region: "), err
	}
	for _, v := range list.Items {
		bt, err := v.MarshalJSON()
        if err != nil {
			return "", err
		}
        rw := json.RawMessage(bt)
        json, _ := json.MarshalIndent(rw, "", "\t")
        res = res + fmt.Sprintf("%s", json)
	}

    return res, nil
}

/* @see https://cloud.google.com/compute/docs/reference/rest/v1/instances/list
 */
func (n *Compute) getInstancesList() (string, error) {
    var res string
    var svc = n.client

    list, err := svc.Instances.List(n.Project, n.Zone).Do()
	if err != nil {
		return fmt.Sprintf("failed to list instances: "), err
	}
	for _, v := range list.Items {
		bt, err := v.MarshalJSON()
        if err != nil {
			return "", err
		}
        rw := json.RawMessage(bt)
        json, _ := json.MarshalIndent(rw, "", "\t")
        res = res + fmt.Sprintf("%s", json)
	}

    return res, nil
}

/*
 */
func (n *Compute) Close() { }

