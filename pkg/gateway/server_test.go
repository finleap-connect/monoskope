package gateway

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ = Describe("Gateway", func() {
	It("declines invalid bearer token", func() {
		conn, err := CreateGatewayConnecton(gatewayApiListener.Addr().String(), env.GatewayClientTransportCredentials, invalidToken())
		if err != nil {
			log.Error(err, "did not connect: %v")
		}
		defer conn.Close()
		gwc := api.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).To(HaveOccurred())
		Expect(serverInfo).To(BeNil())
	})
	It("accepts root bearer token", func() {
		conn, err := CreateGatewayConnecton(gatewayApiListener.Addr().String(), env.GatewayClientTransportCredentials, rootToken())
		if err != nil {
			log.Error(err, "did not connect: %v")
		}
		defer conn.Close()
		gwc := api.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		Expect(serverInfo).ToNot(BeNil())
	})
	It("can go through oidc-flow with existing user", func() {
		handler, err := getClientAuthHandler(env.DexWebEndpoint, redirectURL)
		Expect(err).ToNot(HaveOccurred())
		authCodeURL, err := getAuthURL(handler)
		Expect(err).ToNot(HaveOccurred())

		res, err := httpClient.Get(authCodeURL)
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

		authToken, err := handler.Exchange(context.Background(), authCode)
		Expect(err).ToNot(HaveOccurred())

		conn, err := CreateGatewayConnecton(gatewayApiListener.Addr().String(), env.GatewayClientTransportCredentials, authToken)
		if err != nil {
			log.Error(err, "did not connect: %v")
		}
		defer conn.Close()
		gwc := api.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		Expect(serverInfo).ToNot(BeNil())
	})
})
