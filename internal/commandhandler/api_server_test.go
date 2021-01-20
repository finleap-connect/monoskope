package commandhandler

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	ctx = context.Background()
)

var _ = Describe("HealthCheck", func() {
	It("can do health checks", func() {
		conn, err := grpc.
			NewGrpcConnectionFactory(testEnv.GetApiAddr()).
			WithInsecure().
			WithBlock().
			Build(ctx)
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()

		hc := healthpb.NewHealthClient(conn)
		res, err := hc.Check(ctx, &healthpb.HealthCheckRequest{})
		Expect(err).ToNot(HaveOccurred())
		Expect(res.GetStatus()).To(Equal(healthpb.HealthCheckResponse_SERVING))
	})
})
