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
		conn, err := CreateGatewayConnecton(gatewayApiListener.Addr().String(), invalidToken())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := gateway.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).To(HaveOccurred())
		Expect(serverInfo).To(BeNil())
	})
	It("accepts root bearer token", func() {
		conn, err := CreateGatewayConnecton(gatewayApiListener.Addr().String(), rootToken())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := gateway.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		Expect(serverInfo).ToNot(BeNil())
	})
	It("can retrieve auth url", func() {
		conn, err := CreateGatewayConnecton(gatewayApiListener.Addr().String(), nil)
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := gateway.NewGatewayClient(conn)

		authInfo, err := gwc.GetAuthInformation(context.Background(), &auth.AuthState{CallbackURL: env.AuthConfig.RedirectURI})
		Expect(err).ToNot(HaveOccurred())
		Expect(authInfo).ToNot(BeNil())
		log.Info("AuthCodeURL: " + authInfo.AuthCodeURL)
	})
	It("can go through oidc-flow with existing user", func() {
		conn, err := CreateGatewayConnecton(gatewayApiListener.Addr().String(), nil)
		Expect(err).ToNot(HaveOccurred())
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
		Expect(userInfo).ToNot(BeNil())
		Expect(userInfo.GetEmail()).To(Equal("admin@example.com"))
		log.Info("Received user info", "AccessToken", userInfo.GetAccessToken(), "Expiry", userInfo.GetExpiry().AsTime())

		conn, err = CreateGatewayConnecton(gatewayApiListener.Addr().String(), toToken(userInfo.GetAccessToken()))
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc = gateway.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		Expect(serverInfo).ToNot(BeNil())
	})
	It("fail to go through oidc-flow for non-existing user", func() {
		conn, err := CreateGatewayConnecton(gatewayApiListener.Addr().String(), nil)
		Expect(err).ToNot(HaveOccurred())
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
