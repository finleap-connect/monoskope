package gateway

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api_gw "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	api_gw_auth "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway/auth"
	"golang.org/x/sync/errgroup"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ctx = context.Background()
)

var _ = Describe("Gateway", func() {
	It("declines invalid bearer token", func() {
		conn, err := CreateInsecureGatewayConnecton(ctx, apiListener.Addr().String(), invalidToken())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := api_gw.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).To(HaveOccurred())
		Expect(serverInfo).To(BeNil())
	})
	It("accepts root bearer token", func() {
		conn, err := CreateInsecureGatewayConnecton(ctx, apiListener.Addr().String(), rootToken())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := api_gw.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		Expect(serverInfo).ToNot(BeNil())
	})
	It("can retrieve auth url", func() {
		conn, err := CreateInsecureGatewayConnecton(ctx, apiListener.Addr().String(), nil)
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := api_gw_auth.NewAuthClient(conn)

		authInfo, err := gwc.GetAuthInformation(context.Background(), &api_gw_auth.AuthState{CallbackURL: "http://localhost:8000"})
		Expect(err).ToNot(HaveOccurred())
		Expect(authInfo).ToNot(BeNil())
		env.Log.Info("AuthCodeURL: " + authInfo.AuthCodeURL)
	})
	It("can go through oidc-flow with existing user", func() {
		conn, err := CreateInsecureGatewayConnecton(ctx, apiListener.Addr().String(), nil)
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwcAuth := api_gw_auth.NewAuthClient(conn)

		ready := make(chan string, 1)
		oidcClientServer, err := env.NewOidcClientServer(ready)
		Expect(err).ToNot(HaveOccurred())
		defer oidcClientServer.Close()

		env.Log.Info("oidc redirect uri: " + oidcClientServer.RedirectURI)
		authInfo, err := gwcAuth.GetAuthInformation(context.Background(), &api_gw_auth.AuthState{CallbackURL: oidcClientServer.RedirectURI})
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
			env.Log.Info("wait for oidc client server to get ready...")
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

		authResponse, err := gwcAuth.ExchangeAuthCode(context.Background(), &api_gw_auth.AuthCode{Code: authCode, State: authInfo.GetState(), CallbackURL: oidcClientServer.RedirectURI})
		Expect(err).ToNot(HaveOccurred())
		Expect(authResponse).ToNot(BeNil())
		Expect(authResponse.GetEmail()).To(Equal("admin@example.com"))
		Expect(authResponse.GetAccessToken()).ToNot(Equal(""))
		Expect(authResponse.GetRefreshToken()).ToNot(Equal(""))
		env.Log.Info("Received user info", "AccessToken", authResponse.GetAccessToken(), "Expiry", authResponse.GetAccessToken().GetExpiry().AsTime())

		conn, err = CreateInsecureGatewayConnecton(ctx, apiListener.Addr().String(), toToken(authResponse.GetAccessToken().GetToken()))
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := api_gw.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		Expect(serverInfo).ToNot(BeNil())

		accessToken, err := gwcAuth.RefreshAuth(ctx, &api_gw_auth.RefreshAuthRequest{RefreshToken: authResponse.GetRefreshToken()})

		conn, err = CreateInsecureGatewayConnecton(ctx, apiListener.Addr().String(), toToken(accessToken.GetToken()))
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc = api_gw.NewGatewayClient(conn)

		serverInfo, err = gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		Expect(serverInfo).ToNot(BeNil())

	})
})

var _ = Describe("HealthCheck", func() {
	It("can do health checks", func() {
		conn, err := CreateInsecureGatewayConnecton(ctx, apiListener.Addr().String(), nil)
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()

		hc := healthpb.NewHealthClient(conn)
		res, err := hc.Check(ctx, &healthpb.HealthCheckRequest{})
		Expect(err).ToNot(HaveOccurred())
		Expect(res.GetStatus()).To(Equal(healthpb.HealthCheckResponse_SERVING))
	})
})
