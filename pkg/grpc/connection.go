package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type tlsConf struct {
	pemServerCAFile string
	certFile        string
	keyFile         string
}

type grpcConnectionFactory struct {
	opts    []grpc.DialOption
	url     string
	tlsConf *tlsConf
}

func NewGrpcConnectionFactory(url string) grpcConnectionFactory {
	return grpcConnectionFactory{
		url: url,
	}
}

func NewGrpcConnectionFactoryWithDefaults(url string) grpcConnectionFactory {
	return NewGrpcConnectionFactory(url).
		WithInsecure().
		WithBlock()
}

func (factory grpcConnectionFactory) WithInsecure() grpcConnectionFactory {
	if factory.opts == nil {
		factory.opts = make([]grpc.DialOption, 0)
	}
	factory.opts = append(factory.opts, grpc.WithInsecure())
	return factory
}

func (factory grpcConnectionFactory) WithTransportCredentials(pemServerCAFile, certFile, keyFile string) grpcConnectionFactory {
	factory.tlsConf = &tlsConf{
		pemServerCAFile: pemServerCAFile,
		certFile:        certFile,
		keyFile:         keyFile,
	}
	return factory
}

func (factory grpcConnectionFactory) WithBlock() grpcConnectionFactory {
	if factory.opts == nil {
		factory.opts = make([]grpc.DialOption, 0)
	}
	factory.opts = append(factory.opts, grpc.WithBlock())
	return factory
}

func (factory grpcConnectionFactory) WithRetry() grpcConnectionFactory {
	if factory.opts == nil {
		factory.opts = make([]grpc.DialOption, 0)
	}

	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(10 * time.Millisecond)),
		grpc_retry.WithCodes(grpc_retry.DefaultRetriableCodes...),
		grpc_retry.WithMax(5),
	}

	factory.opts = append(factory.opts, grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)))
	factory.opts = append(factory.opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))

	return factory
}

func (factory grpcConnectionFactory) Connect(ctx context.Context) (*grpc.ClientConn, error) {
	if factory.tlsConf != nil {
		tlsCredentials, err := factory.loadTLSCredentials()
		if err != nil {
			return nil, err
		}
		factory.opts = append(factory.opts, grpc.WithTransportCredentials(tlsCredentials))
	}

	return grpc.DialContext(ctx, factory.url, factory.opts...)
}

func (factory grpcConnectionFactory) ConnectWithTimeout(ctx context.Context, timeout time.Duration) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return factory.Connect(ctx)
}

// loadTLSCredentials actually loads the configured certs
func (factory *grpcConnectionFactory) loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile(factory.tlsConf.pemServerCAFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(factory.tlsConf.certFile, factory.tlsConf.keyFile)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}
