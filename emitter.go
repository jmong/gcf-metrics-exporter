package metricsexporter

/*
 */
type Emitters interface {
    Emit() error
}

/* @TODO
 * Various types of available emitters.
 */
type EmitterType int
const (
    EMITTER_PROMETHEUS = iota
    //EMITTER_STACKDRIVER
)
