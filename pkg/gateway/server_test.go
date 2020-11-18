package gateway

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/api/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/api/gateway/auth"
	"golang.org/x/sync/errgroup"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ctx = context.Background()
)

var _ = Describe("Gateway", func() {
	It("declines invalid bearer token", func() {
		conn, err := CreateGatewayConnecton(ctx, gatewayApiListener.Addr().String(), invalidToken())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := gateway.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).To(HaveOccurred())
		Expect(serverInfo).To(BeNil())
	})
	It("accepts root bearer token", func() {
		conn, err := CreateGatewayConnecton(ctx, gatewayApiListener.Addr().String(), rootToken())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := gateway.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		Expect(serverInfo).ToNot(BeNil())
	})
	It("can retrieve auth url", func() {
		conn, err := CreateGatewayConnecton(ctx, gatewayApiListener.Addr().String(), nil)
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := gateway.NewGatewayClient(conn)

		authInfo, err := gwc.GetAuthInformation(context.Background(), &auth.AuthState{CallbackURL: "http://localhost:8000"})
		Expect(err).ToNot(HaveOccurred())
		Expect(authInfo).ToNot(BeNil())
		log.Info("AuthCodeURL: " + authInfo.AuthCodeURL)
	})
	It("can go through oidc-flow with existing user", func() {
		conn, err := CreateGatewayConnecton(ctx, gatewayApiListener.Addr().String(), nil)
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := gateway.NewGatewayClient(conn)

		ready := make(chan string, 1)
		oidcClientServer, err := env.NewOidcClientServer(ready)
		Expect(err).ToNot(HaveOccurred())
		defer oidcClientServer.Close()

		log.Info("oidc redirect uri: " + oidcClientServer.RedirectURI)
		authInfo, err := gwc.GetAuthInformation(context.Background(), &auth.AuthState{CallbackURL: oidcClientServer.RedirectURI})
		Expect(err).ToNot(HaveOccurred())
		Expect(authInfo).ToNot(BeNil())

		var innerErr error
		res, err := httpClient.Get(authInfo.AuthCodeURL)
		Expect(err).NotTo(HaveOccurred())
		doc, err := goquery.NewDocumentFromReader(res.Body)
		Expect(err).NotTo(HaveOccurred())
		path, ok := doc.Find("form").Attr("action")
		Expect(ok).To(BeTrue())
		formAction := fmt.Sprintf("%s%s", env.DexWebEndpoint, path)

		var authCode string
		var statusCode int
		var eg errgroup.Group
		eg.Go(func() error {
			defer GinkgoRecover()
			var innerErr error
			authCode, innerErr = oidcClientServer.ReceiveCodeViaLocalServer(ctx, authInfo.AuthCodeURL, authInfo.State)
			return innerErr
		})
		eg.Go(func() error {
			defer GinkgoRecover()
			log.Info("wait for oidc client server to get ready...")
			<-ready
			res, err = httpClient.PostForm(formAction, url.Values{
				"login": {"admin@example.com"}, "password": {"password"},
			})
			if err == nil {
				statusCode = res.StatusCode
			}
			return innerErr
		})
		Expect(eg.Wait()).NotTo(HaveOccurred())
		Expect(statusCode).To(Equal(http.StatusOK))

		userInfo, err := gwc.ExchangeAuthCode(context.Background(), &auth.AuthCode{Code: authCode, State: authInfo.GetState(), CallbackURL: oidcClientServer.RedirectURI})
		Expect(err).ToNot(HaveOccurred())
		Expect(userInfo).ToNot(BeNil())
		Expect(userInfo.GetEmail()).To(Equal("admin@example.com"))
		log.Info("Received user info", "AccessToken", userInfo.GetAccessToken(), "Expiry", userInfo.GetExpiry().AsTime())

		conn, err = CreateGatewayConnecton(ctx, gatewayApiListener.Addr().String(), toToken(userInfo.GetAccessToken()))
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc = gateway.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		Expect(serverInfo).ToNot(BeNil())
	})
})

var _ = Describe("HealthCheck", func() {
	It("can do health checks", func() {
		conn, err := CreateGatewayConnecton(ctx, gatewayApiListener.Addr().String(), nil)
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()

		hc := healthpb.NewHealthClient(conn)
		res, err := hc.Check(ctx, &healthpb.HealthCheckRequest{})
		Expect(err).ToNot(HaveOccurred())
		Expect(res.GetStatus()).To(Equal(healthpb.HealthCheckResponse_SERVING))
	})
})
