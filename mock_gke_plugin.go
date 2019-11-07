package metricsexporter

import ()

/*
 *
 */
func NewMockGKEBuilder() GKEBuilder {
    return &gkeBuild{}
}

/*
 *
 */
func (b *gkeBuild) BuildMock() (GKE, error) {
    return GKE{
        context:  b.context,
        client:   nil,
        Cluster:  b.cluster,
        Project:  b.project,
        Zone:     b.zone,
    }, nil
}
