package gateway

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	ctx = context.Background()
)

var _ = Describe("Gateway", func() {
	It("can retrieve auth url", func() {
		conn, err := CreateInsecureGatewayConnection(ctx, apiListener.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := api.NewGatewayClient(conn)

		authInfo, err := gwc.GetAuthInformation(context.Background(), &api.AuthState{CallbackURL: "http://localhost:8000"})
		Expect(err).ToNot(HaveOccurred())
		Expect(authInfo).ToNot(BeNil())
		env.Log.Info("AuthCodeURL: " + authInfo.AuthCodeURL)
	})
})

var _ = Describe("HealthCheck", func() {
	It("can do health checks", func() {
		conn, err := CreateInsecureGatewayConnection(ctx, apiListener.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()

		hc := healthpb.NewHealthClient(conn)
		res, err := hc.Check(ctx, &healthpb.HealthCheckRequest{})
		Expect(err).ToNot(HaveOccurred())
		Expect(res.GetStatus()).To(Equal(healthpb.HealthCheckResponse_SERVING))
	})
})
