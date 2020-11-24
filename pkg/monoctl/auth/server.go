package auth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/int128/listener"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	listener    *listener.Listener
	RedirectURI string
	config      *Config
	log         logger.Logger
}

func NewServer(c *Config) (*Server, error) {
	l, err := listener.New(c.LocalServerBindAddress)
	if err != nil {
		return nil, fmt.Errorf("could not start a local server: %w", err)
	}
	return &Server{
		listener:    l,
		config:      c,
		RedirectURI: computeRedirectURL(l, c),
		log:         logger.WithName("oidc-client-server"),
	}, nil
}

func (s *Server) Close() {
	s.log.Info("Server stopped")
	defer s.listener.Close()
}

func (s *Server) ReceiveCodeViaLocalServer(ctx context.Context, authCodeURL, state string) (string, error) {
	respCh := make(chan *authorizationResponse)
	server := http.Server{
		Handler: &localServerHandler{
			localServerSuccessHTML: s.config.LocalServerSuccessHTML,
			authCodeUrl:            authCodeURL,
			state:                  state,
			respCh:                 respCh,
		},
	}

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
		case s.config.LocalServerReadyChan <- s.RedirectURI:
			s.log.Info("Server ready")
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

type localServerHandler struct {
	state                  string
	authCodeUrl            string
	respCh                 chan<- *authorizationResponse // channel to send a response to
	onceRespCh             sync.Once                     // ensure send once
	localServerSuccessHTML string
}

func (h *localServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	switch {
	case r.Method == "GET" && r.URL.Path == "/" && q.Get("error") != "":
		h.onceRespCh.Do(func() {
			h.respCh <- h.handleErrorResponse(w, r)
		})
	case r.Method == "GET" && r.URL.Path == "/" && q.Get("code") != "":
		h.onceRespCh.Do(func() {
			h.respCh <- h.handleCodeResponse(w, r)
		})
	case r.Method == "GET" && r.URL.Path == "/":
		h.handleIndex(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *localServerHandler) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, h.authCodeUrl, http.StatusFound)
}

func (h *localServerHandler) handleCodeResponse(w http.ResponseWriter, r *http.Request) *authorizationResponse {
	q := r.URL.Query()
	code, state := q.Get("code"), q.Get("state")

	if state != h.state {
		http.Error(w, "authorization error", 500)
		return &authorizationResponse{err: fmt.Errorf("state does not match (wants %s but got %s)", h.state, state)}
	}
	w.Header().Add("Content-Type", "text/html")
	if _, err := fmt.Fprintf(w, h.localServerSuccessHTML); err != nil {
		http.Error(w, "server error", 500)
		return &authorizationResponse{err: fmt.Errorf("write error: %w", err)}
	}
	return &authorizationResponse{code: code}
}

func (h *localServerHandler) handleErrorResponse(w http.ResponseWriter, r *http.Request) *authorizationResponse {
	q := r.URL.Query()
	errorCode, errorDescription := q.Get("error"), q.Get("error_description")

	http.Error(w, "authorization error", 500)
	return &authorizationResponse{err: fmt.Errorf("authorization error from server: %s %s", errorCode, errorDescription)}
}
