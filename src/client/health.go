package client

import (
	"github.com/pachyderm/pachyderm/v2/src/health"
	"github.com/pachyderm/pachyderm/v2/src/internal/errors"
)

// Health health checks pachd, it returns an error if pachd isn't healthy.
func (c APIClient) Health() error {
	var req health.HealthCheckRequest

	response, err := c.healthClient.Check(c.Ctx(), &req)
	if err != nil {
		return errors.Errorf("health check errored %w", err)
	}
	if response.Status == health.HealthCheckResponse_NOT_SERVING {
		return errors.Errorf("server not ready")
	}
	return nil
}
