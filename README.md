# gcf-metrics-exporter

*WORK-IN-PROGRESS*

## Usage

```
# Deploy
gcloud beta functions deploy RunMetricsExporterHttp --source=. --trigger-http --runtime=go111

# Sample GKE query
# You must have a Service Account (with "Kubernetes Engine Cluster Viewer" and "Kubernetes Engine Viewer" roles) attached to the cluster.
gcloud beta functions call RunMetricsExporterHttp --data '{"resource":"gke", "namespace": "my-gke-cluster", "action": "get", "target": "services.list", "project": "my-gcp-project", "zone": "us-central1-a"}'

# Sample Network query
gcloud beta functions call RunMetricsExporterHttp --data '{"resource":"network", "namespace": "foo", "action": "get", "target": "subnets.list", "project": "my-gcp-project", "region": "us-central1"}'
```

### Required Request Fields

Every query requires the following fields in your json request:
* `resource` - The GCP resource<br>
Valid values:
  * **gke**
  * **network**
  * **compute**
  * **health** _(NOT a GCP resource, see below)_
* `action` - What action to take<br>
Valid values:
  * **get** - Fetches information about the resource
  * **ping** - "Healthcheck" signal to this application (only available in _health_ resource)
* `project` - GCP project id where the resource resides in

### Per-Resource Request Fields

Other json request fields are available depending on the resource.

#### For `gke` resource

You must also include the following fields:
* `namespace` - Name of the Kubernetes cluster
* `zone` - GCP zone the Kubernetes cluster resides in
* `target` - Information about the resource you are looking for<br>
Here are currently available values (subject to change):
  * **services.list** - See for details ... https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.zones.clusters/list
  * **nodepools.list** - See for details ... https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.zones.clusters.nodePools/list
  * **usablesubnets.list** - See for details ... https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1beta1/projects.aggregated.usableSubnetworks/list

#### For `network` resource

You must also include the following fields:
* `region` - GCP region of the virtual networking infrastructure
* `target` - Information about the resource you are looking for<br>
Here are currently available values (subject to change):
  * **subnets.list** - See for details ... https://cloud.google.com/compute/docs/reference/rest/v1/subnetworks/list
  * **firewalls.list** - See for details ... https://cloud.google.com/compute/docs/reference/rest/v1/firewalls/list
  * **addresses.list** - See for details ... https://cloud.google.com/compute/docs/reference/rest/v1/addresses/list
  * **globaladdresses.list** - See for details ... https://cloud.google.com/compute/docs/reference/rest/v1/globalAddresses/list
  * **networks.list** - See for details ... https://cloud.google.com/compute/docs/reference/rest/v1/networks/list
  * **routers.list** - See for details ... https://cloud.google.com/compute/docs/reference/rest/v1/routers/list
  * **routes.list** - See for details ... https://cloud.google.com/compute/docs/reference/rest/v1/routes/list
  * **interconnects.list** - See for details ... https://cloud.google.com/compute/docs/reference/rest/v1/interconnects/list

#### For `compute` resource

You must also include the following fields:
* `region` - GCP region where the compute resource resides in
* `target` - Information about the resource you are looking for<br>
Here are currently available values (subject to change):
  * **regions.list** - See for details ... https://cloud.google.com/compute/docs/reference/rest/v1/regions/list
  * **instances.list** - See for details ... https://cloud.google.com/compute/docs/reference/rest/v1/instances/list

#### For `health` resource

This is _not_ a Google Cloud resource. It is used to retrieve certain information about the "health" of this application.

To check "aliveness" of this application, just set the value of the json field "action" to "ping".
```
{"action": "ping"}
```
