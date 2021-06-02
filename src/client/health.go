package client

import (
	"github.com/pachyderm/pachyderm/v2/src/health"
	"github.com/pachyderm/pachyderm/v2/src/internal/grpcutil"
)

// Health health checks pachd, it returns an error if pachd isn't healthy.
func (c APIClient) Health() error {
	var health health.HealthCheckRequest

	_, err := c.healthClient.Check(c.Ctx(), &health)
	return grpcutil.ScrubGRPC(err)
}
