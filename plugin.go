package metricsexporter

/*
 */
type Plugins interface {
    Do(string)  (string, error)
}

/* Common GCP metadata we require for all projects.
 */
type GcpMetadata struct {
    project      string
    environment  string
    timer        Timer
}

//var(
// TBD
//    ConnPool  map[string]*Plugins
//)
