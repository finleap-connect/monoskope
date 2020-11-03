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
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ = Describe("Gateway", func() {
	It("declines invalid bearer token", func() {
		conn, err := CreateGatewayAuthedConnecton(gatewayApiListener.Addr().String(), env.GatewayClientTransportCredentials, invalidToken())
		if err != nil {
			log.Error(err, "did not connect: %v")
		}
		defer conn.Close()
		gwc := gateway.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).To(HaveOccurred())
		Expect(serverInfo).To(BeNil())
	})
	It("accepts root bearer token", func() {
		conn, err := CreateGatewayAuthedConnecton(gatewayApiListener.Addr().String(), env.GatewayClientTransportCredentials, rootToken())
		if err != nil {
			log.Error(err, "did not connect: %v")
		}
		defer conn.Close()
		gwc := gateway.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		Expect(serverInfo).ToNot(BeNil())
	})
	It("can retrieve auth url", func() {
		conn, err := CreateGatewayConnecton(gatewayApiListener.Addr().String(), env.GatewayClientTransportCredentials)
		if err != nil {
			log.Error(err, "did not connect: %v")
		}
		defer conn.Close()
		gwc := gateway.NewGatewayClient(conn)

		authInfo, err := gwc.GetAuthInformation(context.Background(), &auth.AuthState{CallbackURL: env.AuthConfig.RedirectURI})
		Expect(err).ToNot(HaveOccurred())
		Expect(authInfo).ToNot(BeNil())
		log.Info("AuthCodeURL: " + authInfo.AuthCodeURL)
	})
	It("can go through oidc-flow with existing user", func() {
		conn, err := CreateGatewayConnecton(gatewayApiListener.Addr().String(), env.GatewayClientTransportCredentials)
		if err != nil {
			log.Error(err, "did not connect: %v")
		}
		defer conn.Close()
		gwc := gateway.NewGatewayClient(conn)

		authInfo, err := gwc.GetAuthInformation(context.Background(), &auth.AuthState{CallbackURL: env.AuthConfig.RedirectURI})
		Expect(err).ToNot(HaveOccurred())
		Expect(authInfo).ToNot(BeNil())

		res, err := httpClient.Get(authInfo.AuthCodeURL)
		Expect(err).NotTo(HaveOccurred())
		doc, err := goquery.NewDocumentFromReader(res.Body)
		Expect(err).NotTo(HaveOccurred())
		path, ok := doc.Find("form").Attr("action")
		Expect(ok).To(BeTrue())
		res, err = httpClient.PostForm(fmt.Sprintf("%s%s", env.DexWebEndpoint, path), url.Values{
			"login": {"admin@example.com"}, "password": {"password"},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))

		userInfo, err := gwc.ExchangeAuthCode(context.Background(), &auth.AuthCode{Code: authCode, State: authState})
		Expect(err).ToNot(HaveOccurred())

		conn, err = CreateGatewayAuthedConnecton(gatewayApiListener.Addr().String(), env.GatewayClientTransportCredentials, toToken(userInfo.GetAccessToken()))
		if err != nil {
			log.Error(err, "did not connect: %v")
		}
		defer conn.Close()
		gwc = gateway.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		Expect(serverInfo).ToNot(BeNil())
	})
	It("can go through oidc-flow declining non-existing user", func() {
		conn, err := CreateGatewayConnecton(gatewayApiListener.Addr().String(), env.GatewayClientTransportCredentials)
		if err != nil {
			log.Error(err, "did not connect: %v")
		}
		defer conn.Close()
		gwc := gateway.NewGatewayClient(conn)

		authInfo, err := gwc.GetAuthInformation(context.Background(), &auth.AuthState{CallbackURL: env.AuthConfig.RedirectURI})
		Expect(err).ToNot(HaveOccurred())
		Expect(authInfo).ToNot(BeNil())

		res, err := httpClient.Get(authInfo.AuthCodeURL)
		Expect(err).NotTo(HaveOccurred())
		doc, err := goquery.NewDocumentFromReader(res.Body)
		Expect(err).NotTo(HaveOccurred())
		path, ok := doc.Find("form").Attr("action")
		Expect(ok).To(BeTrue())
		res, err = httpClient.PostForm(fmt.Sprintf("%s%s", env.DexWebEndpoint, path), url.Values{
			"login": {"wronguser"}, "password": {"wrongpassword"},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))

		userInfo, err := gwc.ExchangeAuthCode(context.Background(), &auth.AuthCode{Code: authCode, State: authState})
		Expect(err).To(HaveOccurred())
		Expect(userInfo).To((BeNil()))
	})
})
