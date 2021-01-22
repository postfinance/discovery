// Package server is responsible for starting the grpc and http server.
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/postfinance/store"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zbindenren/discovery"
	"github.com/zbindenren/discovery/internal/registry"
	discoveryv1 "github.com/zbindenren/discovery/pkg/discoverypb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	httpStopTimeout = 10 * time.Second
)

// Server represents the discovery server.
type Server struct {
	wg          *sync.WaitGroup
	numReplicas int
	backend     store.Backend
	reg         prometheus.Registerer
	l           *zap.SugaredLogger
	grpc        grpcConfig
	http        httpConfig
}

// New initializes a new Server.
func New(backend store.Backend, l *zap.SugaredLogger, reg prometheus.Registerer, numReplicas int, grpcListen, httpListen string) (*Server, error) {
	s := Server{
		numReplicas: numReplicas,
		backend:     backend,
		reg:         reg,
		l:           l,
		wg:          &sync.WaitGroup{},
		grpc: grpcConfig{
			listenAddr: grpcListen,
		},
		http: httpConfig{
			listenAddr: httpListen,
		},
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

	s.grpc.server = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(opts...),
			grpcMetrics.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(opts...),
			grpcMetrics.UnaryServerInterceptor(),
		)),
	)

	if err := s.reg.Register(grpcMetrics); err != nil {
		return err
	}

	r, err := registry.New(s.backend, s.reg, s.l, s.numReplicas)
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

	s.l.Infow("recalculating services")

	if _, err := r.ReRegisterAllServices(); err != nil {
		return err
	}

	a := &API{
		r: r,
	}

	discoveryv1.RegisterServerAPIServer(s.grpc.server, a)
	discoveryv1.RegisterServiceAPIServer(s.grpc.server, a)
	discoveryv1.RegisterNamespaceAPIServer(s.grpc.server, a)

	listener, err := net.Listen("tcp", s.grpc.listenAddr)

	if err != nil {
		return err
	}

	if err := s.grpc.server.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (s *Server) startHTTP(ctx context.Context) error {
	s.l.Infow("starting http server")

	ep, err := getGRPCEndpointFromListenAddr(s.grpc.listenAddr)
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

	r, ok := s.reg.(prometheus.Gatherer)
	if !ok {
		panic("interface is not prometheus.Registry")
	}

	mux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	mux.Handle("/", gwmux)

	s.http.server = &http.Server{
		Addr:    s.http.listenAddr,
		Handler: mux,
	}

	if err := s.http.server.ListenAndServe(); err != http.ErrServerClosed {
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

	if s.grpc.server == nil {
		return
	}

	s.grpc.server.GracefulStop()
}

func (s *Server) stopHTTP() error {
	defer s.wg.Done()
	defer s.l.Info("http server stopped")

	if s.http.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), httpStopTimeout)
	defer cancel()

	return s.http.server.Shutdown(ctx)
}

type grpcConfig struct {
	server     *grpc.Server
	listenAddr string
}

type httpConfig struct {
	server     *http.Server
	listenAddr string
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
