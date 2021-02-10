// Package server is responsible for starting the grpc and http server.
package server

import (
	"context"
	"embed"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/postfinance/discovery"
	"github.com/postfinance/discovery/internal/auth"
	"github.com/postfinance/discovery/internal/registry"
	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb"
	"github.com/postfinance/store"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	httpStopTimeout = 10 * time.Second
)

var (
	//go:embed swagger/*
	static embed.FS //nolint: gochecknoglobals // no other possibility
)

// Server represents the discovery server.
type Server struct {
	wg         *sync.WaitGroup
	backend    store.Backend
	l          *zap.SugaredLogger
	config     Config
	grpcServer *grpc.Server
	httpServer *http.Server
}

// Config configures the discovery server.
type Config struct {
	PrometheusRegistry prometheus.Registerer
	NumReplicas        int
	GRPCListenAddr     string
	HTTPListenAddr     string
	TokenIssuer        string
	TokenSecretKey     string
	OIDCClient         string
	OIDCRoles          []string
	OIDCURL            string
	Transport          http.RoundTripper
}

// New initializes a new Server.
func New(backend store.Backend, l *zap.SugaredLogger, cfg Config) (*Server, error) {
	s := Server{
		backend: backend,
		l:       l,
		wg:      &sync.WaitGroup{},
		config:  cfg,
	}

	return &s, nil
}

// Run starts the server and runs until context is canceled.
func (s *Server) Run(ctx context.Context) error {
	s.wg.Add(2)

	errChan := make(chan error)

	go func() {
		if err := s.startGRPC(); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := s.startHTTP(ctx); err != nil {
			errChan <- err
		}
	}()

	for {
		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			err := s.stop()
			s.wg.Wait()

			return err
		}
	}
}

func (s *Server) startGRPC() error {
	s.l.Infow("starting grpc server")

	grpcMetrics := grpc_prometheus.NewServerMetrics()
	grpcMetrics.EnableHandlingTimeHistogram()

	panicHandler := func(p interface{}) (err error) {
		s.l.Errorw("panic ocured", "trace", string(debug.Stack()))
		return status.Errorf(codes.Unknown, "%v", p)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(panicHandler),
	}

	tokenHandler := auth.NewTokenHandler(s.config.TokenIssuer, s.config.TokenSecretKey)

	verifier, err := auth.NewVerifier(s.config.OIDCURL, s.config.OIDCClient, 10*time.Second, s.config.Transport)
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(opts...),
			grpcMetrics.StreamServerInterceptor(),
			auth.StreamMethodNameInterceptor(),
			grpc_auth.StreamServerInterceptor(auth.Func(verifier, tokenHandler, s.l.Named("auth"))),
			auth.StreamAuthorizeInterceptor(s.config.OIDCRoles...),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(opts...),
			grpcMetrics.UnaryServerInterceptor(),
			auth.UnaryMethodNameInterceptor(),
			grpc_auth.UnaryServerInterceptor(auth.Func(verifier, tokenHandler, s.l.Named("auth"))),
			auth.UnaryAuthorizeInterceptor(s.config.OIDCRoles...),
		)),
	)

	if err := s.config.PrometheusRegistry.Register(grpcMetrics); err != nil {
		return err
	}

	if err := s.config.PrometheusRegistry.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "discovery_replication_factor",
			Help: "A metric with with constant value showing the configured replication factor.",
		},
		func() float64 { return float64(s.config.NumReplicas) },
	)); err != nil {
		return err
	}

	r, err := registry.New(s.backend, s.config.PrometheusRegistry, s.l, s.config.NumReplicas)
	if err != nil {
		return err
	}

	ns, err := r.ListNamespaces()
	if err != nil {
		return err
	}

	dflt := discovery.DefaultNamespace()

	if ns.Index(dflt.Name) < 0 {
		s.l.Infow("creating default namespace", "name", dflt.Name)

		if _, err := r.RegisterNamespace(*dflt); err != nil {
			return err
		}
	}

	a := &API{
		r:            r,
		tokenHandler: tokenHandler,
	}

	discoveryv1.RegisterServerAPIServer(s.grpcServer, a)
	discoveryv1.RegisterServiceAPIServer(s.grpcServer, a)
	discoveryv1.RegisterNamespaceAPIServer(s.grpcServer, a)
	discoveryv1.RegisterTokenAPIServer(s.grpcServer, a)

	// grpc reflection support
	reflection.Register(s.grpcServer)

	listener, err := net.Listen("tcp", s.config.GRPCListenAddr)

	if err != nil {
		return err
	}

	if err := s.grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (s *Server) startHTTP(ctx context.Context) error {
	s.l.Infow("starting http server")

	ep, err := getGRPCEndpointFromListenAddr(s.config.GRPCListenAddr)
	if err != nil {
		return err
	}

	gwmux := runtime.NewServeMux()

	err = discoveryv1.RegisterServiceAPIHandlerFromEndpoint(ctx, gwmux, ep, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return err
	}

	err = discoveryv1.RegisterServerAPIHandlerFromEndpoint(ctx, gwmux, ep, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	r, ok := s.config.PrometheusRegistry.(prometheus.Gatherer)
	if !ok {
		panic("interface is not prometheus.Registry")
	}

	mux.Handle("/swagger/", http.FileServer(http.FS(static)))
	mux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	mux.Handle("/", gwmux)

	s.httpServer = &http.Server{
		Addr:    s.config.HTTPListenAddr,
		Handler: mux,
	}

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) stop() error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		if err := s.stopHTTP(); err != nil {
			errChan <- err
		}
	}()

	go func() {
		s.stopGRPC()
	}()

	return <-errChan
}

func (s *Server) stopGRPC() {
	defer func() {
		if s.wg != nil {
			s.wg.Done()
		}
	}()
	defer s.l.Infow("grpc server stopped")

	if s.grpcServer == nil {
		return
	}

	s.grpcServer.GracefulStop()
}

func (s *Server) stopHTTP() error {
	defer s.wg.Done()
	defer s.l.Info("http server stopped")

	if s.httpServer == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), httpStopTimeout)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

func getGRPCEndpointFromListenAddr(grpcep string) (string, error) {
	tcpaddr, err := net.ResolveTCPAddr("tcp", grpcep)
	if err != nil {
		return "", fmt.Errorf("%s is not a valid listen address: %w", grpcep, err)
	}

	host := "localhost"

	if tcpaddr.IP != nil && !tcpaddr.IP.IsUnspecified() {
		host = tcpaddr.IP.String()
	}

	return net.JoinHostPort(host, fmt.Sprintf("%d", tcpaddr.Port)), nil
}
