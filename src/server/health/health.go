package health

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Server adds the Ready method to health.HealthServer.
type Server interface {
	grpc_health_v1.HealthServer
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
func (h *healthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	//TODO: Implement health checking per service, for now global only

	if !h.ready {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *healthServer) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {

	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

// Ready tells pachd to start responding positively to Health requests. This
// will cause the node to pass its k8s readiness check.
func (h *healthServer) Ready() {
	h.ready = true
}
