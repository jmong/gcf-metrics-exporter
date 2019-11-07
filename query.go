package metricsexporter

/**/
type Query struct {
    Resource   string `json:"resource"`
    Project    string `json:"project"`
    Zone       string `json:"zone"`
    Region     string `json:"region"`
    Action     string `json:"action"`
    Namespace  string `json:"namespace"`
    Target     string `json:"target"`
    Arg1       string `json:"arg1"`
}

const(
    REQUEST_MAX_LEN = 50
    PING_OK         = "ok"
)

var(
    QueryResources  = []string{"gke", "gke_mock", "health", "network", "compute"}
    QueryActions    = []string{"get", "ping"}
)
