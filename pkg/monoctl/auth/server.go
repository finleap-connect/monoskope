package auth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/int128/listener"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	listener    *listener.Listener
	RedirectURL string
	config      *Config
}

func NewServer(c *Config) (*Server, error) {
	l, err := listener.New(c.LocalServerBindAddress)
	if err != nil {
		return nil, fmt.Errorf("could not start a local server: %w", err)
	}
	return &Server{
		listener:    l,
		config:      c,
		RedirectURL: computeRedirectURL(l, c),
	}, nil
}

func (s *Server) Close() {
	defer s.listener.Close()
}

func (s *Server) ReceiveCodeViaLocalServer(ctx context.Context) (string, error) {
	respCh := make(chan *authorizationResponse)
	server := http.Server{}

	shutdownCh := make(chan struct{})
	var resp *authorizationResponse
	var eg errgroup.Group
	eg.Go(func() error {
		defer close(respCh)
		if err := server.Serve(s.listener); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("could not start HTTP server: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		defer close(shutdownCh)
		select {
		case gotResp, ok := <-respCh:
			if ok {
				resp = gotResp
			}
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	eg.Go(func() error {
		if s.config.LocalServerReadyChan == nil {
			return nil
		}
		select {
		case s.config.LocalServerReadyChan <- s.RedirectURL:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	eg.Go(func() error {
		<-shutdownCh
		// Gracefully shutdown the server in the timeout.
		// If the server has not started, Shutdown returns nil and this returns immediately.
		// If Shutdown has failed, force-close the server.
		ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			_ = server.Close()
			return nil
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return "", fmt.Errorf("authorization error: %w", err)
	}
	if resp == nil {
		return "", errors.New("no authorization response")
	}
	return resp.code, resp.err
}

func computeRedirectURL(l net.Listener, c *Config) string {
	hostPort := fmt.Sprintf("%s:%d", c.RedirectURLHostname, l.Addr().(*net.TCPAddr).Port)
	return "http://" + hostPort
}

type authorizationResponse struct {
	code string // non-empty if a valid code is received
	err  error  // non-nil if an error is received or any error occurs
}
