package metricsexporter
/**
 * Cloud Function client that you query to retrieve metrics about GCP resources.
 * 
 * @see https://cloud.google.com/functions/docs/calling/http
 * 
 * @sample
 * switch resource {
 * case "gke":
 *     gke, err := NewGKEBuilder().Context(context.Background()).Project("my-project-123").Cluster("gke-cluster-1234").Build()
 *     resp, err := gke.Do("{'resource': 'gke', 'action': 'get', 'project': 'my-project-123', 'namespace': 'gke-cluster-1234', 'target': 'pods.list'}")
 *     fmt.Fprintf(httpResponseWriter, "%s", resp)
 * }
 **/

import (
    "encoding/json"
    "fmt"
    "html"
    "io"
    "net/http"

    "golang.org/x/net/context"
    "github.com/jmong/validator"
)

var (
    checkLen       validator.StringChainer
    checkResource  validator.StringChainer
    checkAction    validator.StringChainer
)

/*
 */
func initialize() {
    checkLen      = validator.BuildStrChain().IsAlphaNum().IsMaxLen(REQUEST_MAX_LEN)
    checkResource = validator.BuildStrChain().IsAlphaNum().IsInList(QueryResources)
    checkAction   = validator.BuildStrChain().IsAlphaNum().IsInList(QueryActions)
}

/* RunMetricsExporterHttp is the Cloud Function HTTP entry point.
 * It dispatches to the appropriate plugin to handle your json request.
 */
func RunMetricsExporterHttp(w http.ResponseWriter, r *http.Request) {
    initialize()

    var qry Query
    if err := json.NewDecoder(r.Body).Decode(&qry); err != nil {
        fmt.Fprint(w, err)
        return
    }
    debug(w, qry)

    //Validations
    if checkResource.ValidateStr(qry.Resource) == false {
        fmt.Fprintf(w, "[debug] Resource (%s) failed validations\n", qry.Resource)
        return
    }
    if checkAction.ValidateStr(qry.Action) == false {
        fmt.Fprintf(w, "[debug] Action (%s) failed validations\n", qry.Action)
        return
    }
    if checkLen.ValidateStr(qry.Project) == false {
        fmt.Fprintf(w, "[debug] Project (%s) failed validations\n", qry.Project)
        return
    }

    ctx := context.Background()
    dispatch(ctx, w, qry)
}

/* dispath calls the appropriate plugin based on the values
 * set in qry.
 */
func dispatch(ctx context.Context, w io.Writer, qry Query) {
    switch qry.Resource {
    case "health":
        health, err := NewHealthBuilder().Build()
        if err != nil {
            fmt.Fprintf(w, "Error creating Health client: %s\n", err)
            return
        }
        
        resp, err := health.Do(qry)
        if err != nil {
            fmt.Fprintf(w, "Error health.Do(): %s\n", err)
            return
        }
        fmt.Fprintf(w, "%s", resp)
        health.Close()
    case "gke":
        if checkLen.ValidateStr(qry.Namespace) == false {
            fmt.Fprintf(w, "[debug] Namespace (%s) failed validations\n", qry.Namespace)
            return
        }
        if checkLen.ValidateStr(qry.Target) == false {
            fmt.Fprintf(w, "[debug] Target (%s) failed validations\n", qry.Target)
            return
        }
        if checkLen.ValidateStr(qry.Zone) == false {
            fmt.Fprintf(w, "[debug] Zone (%s) failed validations\n", qry.Zone)
            return
        }

        gke, err := NewGKEBuilder().Context(ctx).Project(qry.Project).Zone(qry.Zone).Cluster(qry.Namespace).Arg1(qry.Arg1).Build()
        if err != nil {
            fmt.Fprintf(w, "Error creating GKE client: %s\n", err)
            return
        }

        resp, err := gke.Do(qry)
        if err != nil {
            fmt.Fprintf(w, "Error gke.Do(): %s\n", err)
            return
        }
        fmt.Fprintf(w, "%s", resp)
        gke.Close()
    case "network":
        if checkLen.ValidateStr(qry.Namespace) == false {
            fmt.Fprintf(w, "[debug] Namespace (%s) failed validations\n", qry.Namespace)
            return
        }
        if checkLen.ValidateStr(qry.Target) == false {
            fmt.Fprintf(w, "[debug] Target (%s) failed validations\n", qry.Target)
            return
        }
        if checkLen.ValidateStr(qry.Region) == false {
            fmt.Fprintf(w, "[debug] Region (%s) failed validations\n", qry.Region)
            return
        }

        net, err := NewNetworkBuilder().Context(ctx).Project(qry.Project).Region(qry.Region).Build()
        if err != nil {
            fmt.Fprintf(w, "Error creating Compute client: %s\n", err)
            return
        }

        resp, err := net.Do(qry)
        if err != nil {
            fmt.Fprintf(w, "Error gke.Do(): %s\n", err)
            return
        }
        fmt.Fprintf(w, "%s", resp)
        net.Close()
    case "compute":
        if checkLen.ValidateStr(qry.Namespace) == false {
            fmt.Fprintf(w, "[debug] Namespace (%s) failed validations\n", qry.Namespace)
            return
        }
        if checkLen.ValidateStr(qry.Target) == false {
            fmt.Fprintf(w, "[debug] Target (%s) failed validations\n", qry.Target)
            return
        }
        if checkLen.ValidateStr(qry.Region) == false && checkLen.ValidateStr(qry.Zone) == false {
            fmt.Fprintf(w, "[debug] Both Zone (%s) and Region (%s) failed validations\n", qry.Zone, qry.Region)
            return
        }

        comp, err := NewComputeBuilder().Context(ctx).Project(qry.Project).Region(qry.Region).Zone(qry.Zone).Build()
        if err != nil {
            fmt.Fprintf(w, "Error creating Compute client: %s\n", err)
            return
        }

        resp, err := comp.Do(qry)
        if err != nil {
            fmt.Fprintf(w, "Error gke.Do(): %s\n", err)
            return
        }
        fmt.Fprintf(w, "%s", resp)
        comp.Close()
    // FOR TESTING ONLY //
    case "gke_mock":
        if checkLen.ValidateStr(qry.Namespace) == false {
            fmt.Fprintf(w, "[debug] Namespace (%s) failed validations\n", qry.Namespace)
            return
        }
        if checkLen.ValidateStr(qry.Target) == false {
            fmt.Fprintf(w, "[debug] Target (%s) failed validations\n", qry.Target)
            return
        }
        if checkLen.ValidateStr(qry.Zone) == false {
            fmt.Fprintf(w, "[debug] Zone (%s) failed validations\n", qry.Zone)
            return
        }

        gke, err := NewMockGKEBuilder().Context(ctx).Project(qry.Project).Zone(qry.Zone).Cluster(qry.Namespace).BuildMock()
        if err != nil {
            fmt.Fprintf(w, "Error creating mock GKE client: %s\n", err)
            return
        }

        resp, err := gke.Do(qry)
        if err != nil {
            fmt.Fprintf(w, "Error: %s\n", err)
            return
        }
        fmt.Fprintf(w, "%s", resp)
        gke.Close()
    }
}

/* For debugging only
 */
func debug(w io.Writer, qry Query) {
    fmt.Fprintf(w, "[Debug] Resource = %s\n", html.EscapeString(qry.Resource))
    fmt.Fprintf(w, "[Debug] Project = %s\n", html.EscapeString(qry.Project))
    fmt.Fprintf(w, "[Debug] Action = %s\n", html.EscapeString(qry.Action))
    fmt.Fprintf(w, "[Debug] Namespace = %s\n", html.EscapeString(qry.Namespace))
    fmt.Fprintf(w, "[Debug] Target = %s\n", html.EscapeString(qry.Target))
    fmt.Fprintf(w, "[Debug] Arg1 = %s\n", html.EscapeString(qry.Arg1))
    fmt.Fprintf(w, "[Debug] Zone = %s\n", html.EscapeString(qry.Zone))
    fmt.Fprintf(w, "[Debug] Region = %s\n", html.EscapeString(qry.Region))
}

