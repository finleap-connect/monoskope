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
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth"
	auth_client "gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth/client"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ = Describe("Gateway", func() {
	It("declines invalid bearer token", func() {
		perRPC := oauth.NewOauthAccess(invalidToken())

		opts := []grpc.DialOption{
			// In addition to the following grpc.DialOption, callers may also use
			// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
			// itself.
			// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
			grpc.WithPerRPCCredentials(perRPC),
			// oauth.NewOauthAccess requires the configuration of transport
			// credentials.
			grpc.WithTransportCredentials(clientTransportCredentials),
		}

		opts = append(opts, grpc.WithBlock())
		conn, err := grpc.Dial(apiLis.Addr().String(), opts...)
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
		perRPC := oauth.NewOauthAccess(rootToken())

		opts := []grpc.DialOption{
			// In addition to the following grpc.DialOption, callers may also use
			// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
			// itself.
			// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
			grpc.WithPerRPCCredentials(perRPC),
			// oauth.NewOauthAccess requires the configuration of transport
			// credentials.
			grpc.WithTransportCredentials(clientTransportCredentials),
		}

		opts = append(opts, grpc.WithBlock())
		conn, err := grpc.Dial(apiLis.Addr().String(), opts...)
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
		var state auth.State

		handler, err := auth_client.NewHandler(&auth_client.Config{
			BaseConfig: auth.BaseConfig{
				IssuerURL:      dexWebEndpoint,
				OfflineAsScope: true,
			},
			RedirectURI:  redirectURL,
			Nonce:        "secret-nonce",
			ClientId:     "monoctl",
			ClientSecret: "monoctl-app-secret",
		})
		Expect(err).ToNot(HaveOccurred())

		authCodeURL, err := handler.GetAuthCodeURL(&state, &auth.AuthCodeURLConfig{
			Scopes:        []string{"offline_access"},
			Clients:       []string{},
			OfflineAccess: true,
		})
		Expect(err).ToNot(HaveOccurred())

		res, err := httpClient.Get(authCodeURL)
		Expect(err).NotTo(HaveOccurred())
		doc, err := goquery.NewDocumentFromReader(res.Body)
		Expect(err).NotTo(HaveOccurred())
		path, ok := doc.Find("form").Attr("action")
		Expect(ok).To(BeTrue())
		res, err = httpClient.PostForm(fmt.Sprintf("%s%s", dexWebEndpoint, path), url.Values{
			"login": {"admin@example.com"}, "password": {"password"},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))

		validAuthToken, err := handler.Exchange(context.Background(), authCode)
		Expect(err).ToNot(HaveOccurred())

		perRPC := oauth.NewOauthAccess(validAuthToken)

		opts := []grpc.DialOption{
			// In addition to the following grpc.DialOption, callers may also use
			// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
			// itself.
			// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
			grpc.WithPerRPCCredentials(perRPC),
			// oauth.NewOauthAccess requires the configuration of transport
			// credentials.
			grpc.WithTransportCredentials(clientTransportCredentials),
		}

		opts = append(opts, grpc.WithBlock())
		conn, err := grpc.Dial(apiLis.Addr().String(), opts...)
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

func invalidToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-invalid-token",
	}
}

func rootToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: authRootToken,
	}
}
