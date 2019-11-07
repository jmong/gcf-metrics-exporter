package metricsexporter
/**
 * Client to retrieve information about GCP virtual networking infrastructure.
 * 
 * @see REST call definitions - https://cloud.google.com/compute/docs/reference/rest/v1/
 **/

import (
    "fmt"
    "context"
    "encoding/json"

    compute  "google.golang.org/api/compute/v1"
)

type NetworkBuilder interface {
    Context(context.Context)  NetworkBuilder
    Project(string)           NetworkBuilder
    Region(string)            NetworkBuilder
    EnableEmitter()           NetworkBuilder

    Build()                   (Network, error)
}

type networkBuild struct {
    context        context.Context
    project        string
    region         string
    enableemitter  bool
}

/* Network object.
 */
type Network struct {
    GcpMetadata
    context        context.Context
    client         *compute.Service
    Project        string
    Region         string
    EnableEmitter  bool
    emitter        Emitters
}

/* NewNetworkBuilder creates a builder object by adding components/features 
 * that will create a Network object.
 */
func NewNetworkBuilder() NetworkBuilder {
    return &networkBuild{}
}

/* Context is the Google background context of the request.
 */
func (b *networkBuild) Context(ctx context.Context) NetworkBuilder {
	b.context = ctx
	return b
}

/* Project is the GCP project id.
 */
func (b *networkBuild) Project(project string) NetworkBuilder {
	b.project = project
	return b
}

/* Region is the GCP region the networking resource resides in.
 */
func (b *networkBuild) Region(region string) NetworkBuilder {
	b.region = region
	return b
}

/*
 */
func (b *networkBuild) EnableEmitter() NetworkBuilder {
	b.enableemitter = true
	return b
}

/* Build creates a Network object that retrieves information
 * about the GCP networking infrastructure.
 */
func (b *networkBuild) Build() (Network, error) {
    client, err := compute.NewService(b.context)
	if err != nil {
		return Network{}, err
	}

    var pusher *PrometheusPush = nil
    if b.enableemitter == true {
        pusher = NewPrometheusPush(PROM_PUSHGW_URL, PROM_PUSHGW_JOB)
    }

    return Network{
        context:        b.context,
        client:         client,
        Project:        b.project,
        Region:         b.region,
        EnableEmitter:  b.enableemitter,
        emitter:        pusher,
    }, nil
}

/* @TODO
 * Do acts on your request to retrieve and return a response to you.
 */
func (n *Network) Do(qry Query) (string, error) {
    if qry.Resource == "network" && qry.Action == "get" && qry.Target == "subnets.list" {
        return n.getSubnetsList()
    } else if qry.Resource == "network" && qry.Action == "get" && qry.Target == "firewalls.list" {
        return n.getFirewallsList()
    } else if qry.Resource == "network" && qry.Action == "get" && qry.Target == "addresses.list" {
        return n.getAddressesList()
    } else if qry.Resource == "network" && qry.Action == "get" && qry.Target == "globaladdresses.list" {
        return n.getGlobalAddressesList()
    } else if qry.Resource == "network" && qry.Action == "get" && qry.Target == "networks.list" {
        return n.getNetworksList()
    } else if qry.Resource == "network" && qry.Action == "get" && qry.Target == "routers.list" {
        return n.getRoutersList()
    } else if qry.Resource == "network" && qry.Action == "get" && qry.Target == "routes.list" {
        return n.getRoutesList()
    } else if qry.Resource == "network" && qry.Action == "get" && qry.Target == "interconnects.list" {
        return n.getInterconnectsList()
    }
    return "[Debug] It will call some Compute operations to return json response", nil
}

/* @see https://cloud.google.com/compute/docs/reference/rest/v1/subnetworks/list
 */
func (n *Network) getSubnetsList() (string, error) {
    var res string
    var svc = n.client
    
    list, err := svc.Subnetworks.List(n.Project, n.Region).Do()
	if err != nil {
		return fmt.Sprintf("failed to list subnetworks: "), err
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

/* @see https://cloud.google.com/compute/docs/reference/rest/v1/firewalls/list
 */
func (n *Network) getFirewallsList() (string, error) {
    var res string
    var svc = n.client
    
    list, err := svc.Firewalls.List(n.Project).Do()
	if err != nil {
		return fmt.Sprintf("failed to list firewalls: "), err
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

/* @see https://cloud.google.com/compute/docs/reference/rest/v1/addresses/list
 */
func (n *Network) getAddressesList() (string, error) {
    var res string
    var svc = n.client
    
    list, err := svc.Addresses.List(n.Project, n.Region).Do()
	if err != nil {
		return fmt.Sprintf("failed to list addresses: "), err
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

/* @see https://cloud.google.com/compute/docs/reference/rest/v1/globalAddresses/list
 */
func (n *Network) getGlobalAddressesList() (string, error) {
    var res string
    var svc = n.client
    
    list, err := svc.GlobalAddresses.List(n.Project).Do()
	if err != nil {
		return fmt.Sprintf("failed to list global addresses: "), err
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

/* @see https://cloud.google.com/compute/docs/reference/rest/v1/networks/list
 */
func (n *Network) getNetworksList() (string, error) {
    var res string
    var svc = n.client
    
    list, err := svc.Networks.List(n.Project).Do()
	if err != nil {
		return fmt.Sprintf("failed to list networks: "), err
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

/* @see https://cloud.google.com/compute/docs/reference/rest/v1/routers/list
 */
func (n *Network) getRoutersList() (string, error) {
    var res string
    var svc = n.client
    
    list, err := svc.Routers.List(n.Project, n.Region).Do()
	if err != nil {
		return fmt.Sprintf("failed to list routers: "), err
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

/* @see https://cloud.google.com/compute/docs/reference/rest/v1/routes/list
 */
func (n *Network) getRoutesList() (string, error) {
    var res string
    var svc = n.client
    
    list, err := svc.Routes.List(n.Project).Do()
	if err != nil {
		return fmt.Sprintf("failed to list routes: "), err
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

/* @see https://cloud.google.com/compute/docs/reference/rest/v1/interconnects/list
 */
func (n *Network) getInterconnectsList() (string, error) {
    var res string
    var svc = n.client
    
    list, err := svc.Interconnects.List(n.Project).Do()
	if err != nil {
		return fmt.Sprintf("failed to list interconnects: "), err
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
func (n *Network) Close() { }

