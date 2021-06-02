package health

import (
	"github.com/pachyderm/pachyderm/v2/src/health"
	"golang.org/x/net/context"
)

// Server adds the Ready method to health.HealthServer.
type Server interface {
	health.HealthServer
	Ready()
}

// NewHealthServer returns a new health server
func NewHealthServer() Server {
	return &healthServer{}
}

type healthServer struct {
	ready bool
}

// Health implements the Health method for healthServer.
func (h *healthServer) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	//TODO: Implement health checking per service, for now global only

	if !h.ready {
		return &health.HealthCheckResponse{
			Status: health.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

// Ready tells pachd to start responding positively to Health requests. This
// will cause the node to pass its k8s readiness check.
func (h *healthServer) Ready() {
	h.ready = true
}
