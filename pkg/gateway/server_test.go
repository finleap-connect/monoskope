package gateway

import (
	"context"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/api"
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
			log.Fatalf("did not connect: %v", err)
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
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		gwc := api.NewGatewayClient(conn)

		serverInfo, err := gwc.GetServerInfo(context.Background(), &emptypb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		Expect(serverInfo).ToNot(BeNil())
	})
})

// invalidToken simulates a token lookup and omits the details of proper token
// acquisition. For examples of how to acquire an OAuth2 token, see:
// https://godoc.org/golang.org/x/oauth2
func invalidToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}

func rootToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: authRootToken,
	}
}
