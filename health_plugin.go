package metricsexporter
/**
 * Client to retrieve "health" information about this application.
 **/

import (
    "fmt"
)

/*
 */
type HealthBuilder interface {
    Build()  (Health, error)
}

/*
 */
type healthBuild struct {
    timer  *Timer
}

/*
 */
func NewHealthBuilder() HealthBuilder {
    return &healthBuild{
        timer:  NewTimer(),
    }
}

/*
 */
func (b *healthBuild) Build() (Health, error) {
    return Health{}, nil
}

/*
 */
type Health struct {
    GcpMetadata
}

/* @TODO
 */
func (h *Health) Do(qry Query) (string, error) {
    if qry.Action == "ping" {
        return fmt.Sprintf("%s", PING_OK), nil
    } else if qry.Action == "get" && qry.Target == "stats" {
        return h.getStats()
    }
    return "[Debug] It will call some GKE operations to return json response", nil
}

/* @TODO
 */
func (h *Health) getStats() (string, error) {
    return "Some statistics about this running cloud function instance (eg- version, TBD)", nil
}

/*
 */
func (h *Health) Close() { }

