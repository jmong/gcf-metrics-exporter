package metricsexporter
/**
 * Client to retrieve information about a GKE cluster.
 * 
 * @see REST call definitions - https://cloud.google.com/kubernetes-engine/docs/reference/rest/
 **/

import (
    "fmt"
    "context"
    "encoding/json"

    oauth2  "golang.org/x/oauth2/google"
    gke     "google.golang.org/api/container/v1"
)

/**/
/*
var (
    instance *kubernetes.Client
)

func NewSingleton() *kubernetes.Client {
    lock.Lock()
    defer lock.Unlock()

    if instance == nil {
        //instance = new(kubernetes.Client)
        instance = gke.NewClient(ctx,...)
    }

    return instance
}
*/

type GKEBuilder interface {
    Context(context.Context)  GKEBuilder
    Project(string)           GKEBuilder
    Cluster(string)           GKEBuilder
    Zone(string)              GKEBuilder
    Arg1(string)              GKEBuilder
    EnableEmitter()           GKEBuilder

    Build()                   (GKE, error)
    BuildMock()               (GKE, error)
}

type gkeBuild struct {
    context        context.Context
    project        string
    cluster        string
    zone           string
    arg1           string
    enableemitter  bool
}

/* GKE object.
 */
type GKE struct {
    GcpMetadata
    context        context.Context
    client         *gke.Service
    Project        string
    Cluster        string
    Zone           string
    Arg1           string
    EnableEmitter  bool
    emitter        Emitters
}

/* NewGKEBuilder creates a builder object by adding components/features 
 * that will create a GKE object.
 */
func NewGKEBuilder() GKEBuilder {
    return &gkeBuild{enableemitter: false}
}

/* Context is the Google background context of the request.
 */
func (b *gkeBuild) Context(ctx context.Context) GKEBuilder {
	b.context = ctx
	return b
}

/* Cluster is the name of the Kubernetes cluster.
 */
func (b *gkeBuild) Cluster(cluster string) GKEBuilder {
	b.cluster = cluster
	return b
}

/* Project is the GCP project id.
 */
func (b *gkeBuild) Project(project string) GKEBuilder {
	b.project = project
	return b
}

/* Zone is the GCP zone the resource resides in.
 */
func (b *gkeBuild) Zone(zone string) GKEBuilder {
	b.zone = zone
	return b
}

/* Arg1 is an optional argument passed into the requested resource,
 * action, and target. 
 * The argument type and value depend on the context of the request.
 * For example, if you are requesting ???
 */
func (b *gkeBuild) Arg1(arg1 string) GKEBuilder {
	b.arg1 = arg1
	return b
}

/*  
 */
func (b *gkeBuild) EnableEmitter() GKEBuilder {
	b.enableemitter = true
	return b
}

/* Build creates a GKE object that retrieves information
 * about the GCP Kubernetes cluster.
 */
func (b *gkeBuild) Build() (GKE, error) {
    googclt, err := oauth2.DefaultClient(b.context, gke.CloudPlatformScope)
	if err != nil {
		return GKE{}, err
	}
	client, err := gke.New(googclt)
	if err != nil {
		return GKE{}, err
	}

    var pusher *PrometheusPush = nil
    if b.enableemitter == true {
        pusher = NewPrometheusPush(PROM_PUSHGW_URL, PROM_PUSHGW_JOB)
    }

    return GKE{
        context:        b.context,
        client:         client,
        Cluster:        b.cluster,
        Project:        b.project,
        Zone:           b.zone,
        Arg1:           b.arg1,
        EnableEmitter:  b.enableemitter,
        emitter:        pusher,
    }, nil
}

/* @TODO
 * Do acts on your request to retrieve and return a response to you.
 */
func (g *GKE) Do(qry Query) (string, error) {
    if qry.Resource == "gke" && qry.Action == "get" && qry.Target == "pods.list" {
        return g.getPodsList()
    } else if qry.Resource == "gke" && qry.Action == "get" && qry.Target == "services.list" {
        return g.getServicesList()
    } else if qry.Resource == "gke" && qry.Action == "get" && qry.Target == "nodepools.list" {
        return g.getNodePoolsList()
    } else if qry.Resource == "gke" && qry.Action == "get" && qry.Target == "nodepools.get" {
        return g.getNodePoolsGet()
    }  else if qry.Resource == "gke" && qry.Action == "get" && qry.Target == "usablesubnets.list" {
        return g.getUsableSubnetsList()
    }
    return "[Debug] It will call some GKE operations to return json response", nil
}

/* @TODO
 */
func (g *GKE) getPodsList() (string, error) {
    /*
    pods, err := g.client.GetPods(g.context)
    if err != nil {
        return fmt.Sprintf("Error, getPodsList()"), err
    }
    
    var res string
    for _, pod := range pods {
        res = res + "Name:" + pod.ObjectMeta.Name
    }
    
    return res, nil
    */

    return "[Debug] Feature not implemented yet", nil
}

/* @see https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.zones.clusters/list
 */
func (g *GKE) getServicesList() (string, error) {
    var res string
    var svc = g.client
    
    list, err := svc.Projects.Zones.Clusters.List(g.Project, g.Zone).Do()
	if err != nil {
		return fmt.Sprintf("failed to list clusters: "), err
	}
	for _, v := range list.Clusters {
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

/* @see https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.zones.clusters.nodePools/list
 */
func (g *GKE) getNodePoolsList() (string, error) {
    var res string
    var svc = g.client
    
    list, err := svc.Projects.Zones.Clusters.NodePools.List(g.Project, g.Zone, g.Cluster).Do()
	if err != nil {
		return fmt.Sprintf("failed to list node pools: "), err
	}
	for _, v := range list.NodePools {
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
func (g *GKE) getNodePoolsGet() (string, error) {
    var res string
    var svc = g.client
    
    // @see https://godoc.org/google.golang.org/api/container/v1#ProjectsZonesClustersNodePoolsService.Get
    get, err := svc.Projects.Zones.Clusters.NodePools.Get(g.Project, g.Zone, g.Cluster, g.Arg1).Do()
	if err != nil {
		return fmt.Sprintf("failed to get node pools: "), err
	}
    bt, err := get.MarshalJSON()
    if err != nil {
        return "", err
    }
    rw := json.RawMessage(bt)
    json, _ := json.MarshalIndent(rw, "", "\t")
    res = res + fmt.Sprintf("%s", json)

    return res, nil
}

/* @see https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1beta1/projects.aggregated.usableSubnetworks/list
 */
func (g *GKE) getUsableSubnetsList() (string, error) {
    var res string
    var svc = g.client
    
    // @see https://godoc.org/google.golang.org/api/container/v1#ProjectsAggregatedUsableSubnetworksService.List
    get, err := svc.Projects.Aggregated.UsableSubnetworks.List("projects/" + g.Project).Do()
	if err != nil {
		return fmt.Sprintf("failed to get node in node pools: "), err
	}
    bt, err := get.MarshalJSON()
    if err != nil {
        return "", err
    }
    rw := json.RawMessage(bt)
    json, _ := json.MarshalIndent(rw, "", "\t")
    res = res + fmt.Sprintf("%s", json)

    return res, nil
}


/*
 */
func (g *GKE) Close() { }

