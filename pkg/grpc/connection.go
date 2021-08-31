// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// NewGrpcConnectionFactory creates a new factory for gRPC connections.
func NewGrpcConnectionFactory(url string) grpcConnectionFactory {
	return grpcConnectionFactory{
		url: url,
	}
}

// NewGrpcConnectionFactoryWithDefaults creates a new factory for gRPC connections and adds the following dial options: WithInsecure, WithBlock.
func NewGrpcConnectionFactoryWithDefaults(url string) grpcConnectionFactory {
	return NewGrpcConnectionFactory(url).
		WithInsecure().
		WithBlock()
}

// WithInsecure adds a DialOption which disables transport security for this connection. Note that transport security is required unless WithInsecure is set.
func (factory grpcConnectionFactory) WithInsecure() grpcConnectionFactory {
	if factory.opts == nil {
		factory.opts = make([]grpc.DialOption, 0)
	}
	factory.opts = append(factory.opts, grpc.WithInsecure())
	return factory
}

// WithOSCaTransportCredentials adds a DialOption which configures a connection level security credentials (e.g., TLS/SSL) using the CAs known to the OS.
func (factory grpcConnectionFactory) WithOSCaTransportCredentials() grpcConnectionFactory {
	if factory.opts == nil {
		factory.opts = make([]grpc.DialOption, 0)
	}
	factory.opts = append(factory.opts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	return factory
}

// WithTransportCredentials adds a DialOption which configures a connection level security credentials (e.g., TLS/SSL) using the given certificates.
func (factory grpcConnectionFactory) WithTransportCredentials(pemServerCAFile, certFile, keyFile string) grpcConnectionFactory {
	factory.tlsConf = &tlsConf{
		pemServerCAFile: pemServerCAFile,
		certFile:        certFile,
		keyFile:         keyFile,
	}
	return factory
}

// WithPerRPCCredentials adds a DialOption which sets credentials and places auth state on each outbound RPC.
func (factory grpcConnectionFactory) WithPerRPCCredentials(creds credentials.PerRPCCredentials) grpcConnectionFactory {
	if factory.opts == nil {
		factory.opts = make([]grpc.DialOption, 0)
	}
	factory.opts = append(factory.opts, grpc.WithPerRPCCredentials(creds))
	return factory
}

// WithBlock adds a DialOption which makes caller of Dial blocks until the underlying connection is up. Without this, Dial returns immediately and connecting the server happens in background.
func (factory grpcConnectionFactory) WithBlock() grpcConnectionFactory {
	if factory.opts == nil {
		factory.opts = make([]grpc.DialOption, 0)
	}
	factory.opts = append(factory.opts, grpc.WithBlock())
	return factory
}

// WithRetry adds retrying with exponential backoff using the default retryable codes from grpc_retry.DefaultRetriableCodes.
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

// Connect creates a client connection based on the factory.
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

// ConnectWithTimeout creates a client connection based on the factory with a given timeout.
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
