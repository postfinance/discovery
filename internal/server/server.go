// Package server is responsible for starting the grpc and http server.
package server

import (
	"context"
	"embed"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"connectrpc.com/connect"
	"github.com/postfinance/discovery"
	"github.com/postfinance/discovery/internal/registry"
	discoveryv1connect "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1/discoveryv1connect"
	"github.com/postfinance/store"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.pnet.ch/linux/go/auth/self"
	"go.uber.org/zap"
)

const (
	httpStopTimeout              = 10 * time.Second
	httpReadTimeout              = 10 * time.Second
	httpClientTimeout            = 10 * time.Second
	cacheSyncInterval            = 1 * time.Minute
	serviceCounterUpdateInterval = 15 * time.Second
)

//go:embed swagger/*
var static embed.FS

// Server represents the discovery server.
type Server struct {
	ready      atomic.Bool
	wg         *sync.WaitGroup
	backend    store.Backend
	l          *zap.SugaredLogger
	config     Config
	httpServer *http.Server
}

// Config configures the discovery server.
type Config struct {
	PrometheusRegistry prometheus.Registerer
	NumReplicas        int
	ListenAddr         string
	TokenHandler       *self.TokenHandler
	Interceptors       []connect.Interceptor
}

// New initializes a new Server.
func New(backend store.Backend, l *zap.SugaredLogger, cfg Config) (*Server, error) {
	s := Server{
		backend: backend,
		l:       l,
		wg:      &sync.WaitGroup{},
		config:  cfg,
		ready:   atomic.Bool{},
	}

	s.ready.Store(false)

	return &s, nil
}

// Run starts the server and runs until context is canceled.
func (s *Server) Run(ctx context.Context) error {
	var (
		wg   = new(sync.WaitGroup)
		errC = make(chan error)
		err  error
	)

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := s.startHTTP(ctx); err != nil {
			errC <- err
		}
	}()

	select {
	case <-ctx.Done():
	case e := <-errC:
		err = e
	}

	s.stop()
	wg.Wait()

	return err
}

func (s *Server) createMux(api *API) *http.ServeMux {
	mux := http.NewServeMux()

	r, ok := s.config.PrometheusRegistry.(prometheus.Gatherer)
	if !ok {
		panic("interface is not prometheus.Registry")
	}

	mux.HandleFunc("/ready", s.readyHandler())
	mux.Handle("/swagger/", http.FileServer(http.FS(static)))
	mux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	httpPath, handler := discoveryv1connect.NewServerAPIHandler(api, connect.WithInterceptors(s.config.Interceptors...))
	mux.Handle(httpPath, handler)
	httpPath, handler = discoveryv1connect.NewNamespaceAPIHandler(api, connect.WithInterceptors(s.config.Interceptors...))
	mux.Handle(httpPath, handler)
	httpPath, handler = discoveryv1connect.NewTokenAPIHandler(api, connect.WithInterceptors(s.config.Interceptors...))
	mux.Handle(httpPath, handler)
	httpPath, handler = discoveryv1connect.NewServiceAPIHandler(api, connect.WithInterceptors(s.config.Interceptors...))
	mux.Handle(httpPath, handler)

	return mux
}

func (s *Server) startHTTP(ctx context.Context) error {
	s.l.Infow("starting http server")

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

	go r.StartCacheUpdater(ctx, cacheSyncInterval)
	go r.StartServiceCounterUpdater(ctx, serviceCounterUpdateInterval)

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
		tokenHandler: s.config.TokenHandler,
	}

	mux := s.createMux(a)

	s.httpServer = &http.Server{
		Addr:        s.config.ListenAddr,
		Handler:     mux,
		ReadTimeout: httpClientTimeout,
	}

	s.ready.Store(true)

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

//	func (s *Server) startGRPC(ctx context.Context) error {
//		s.l.Infow("starting grpc server")
//
//		grpcMetrics := grpc_prometheus.NewServerMetrics()
//		grpcMetrics.EnableHandlingTimeHistogram()
//
//		panicHandler := func(p interface{}) (err error) {
//			s.l.Errorw("panic ocured", "trace", string(debug.Stack()))
//			return status.Errorf(codes.Unknown, "%v", p)
//		}
//		opts := []grpc_recovery.Option{
//			grpc_recovery.WithRecoveryHandler(panicHandler),
//		}
//
//		tokenHandler := auth.NewTokenHandler(s.config.TokenIssuer, s.config.TokenSecretKey)
//
//		verifier, err := auth.NewVerifier(s.config.OIDCURL, s.config.OIDCClient, httpClientTimeout, s.config.Transport)
//		if err != nil {
//			return err
//		}
//
//		// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
//		// grpc_zap.ReplaceGrpcLoggerV2(s.l.Desugar())
//
//		s.grpcServer = grpc.NewServer(
//			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
//				grpc_recovery.StreamServerInterceptor(opts...),
//				grpcMetrics.StreamServerInterceptor(),
//				auth.StreamMethodNameInterceptor(),
//				grpc_auth.StreamServerInterceptor(auth.Func(verifier, tokenHandler, s.l.Named("auth"), s.config.ClaimConfig)),
//				auth.StreamAuthorizeInterceptor(s.config.OIDCRoles...),
//				grpc_zap.StreamServerInterceptor(s.l.Desugar(), grpc_zap.WithLevels(customCodeToLevel)),
//			)),
//			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
//				grpc_recovery.UnaryServerInterceptor(opts...),
//				grpcMetrics.UnaryServerInterceptor(),
//				auth.UnaryMethodNameInterceptor(),
//				grpc_auth.UnaryServerInterceptor(auth.Func(verifier, tokenHandler, s.l.Named("auth"), s.config.ClaimConfig)),
//				auth.UnaryAuthorizeInterceptor(s.config.OIDCRoles...),
//				grpc_zap.UnaryServerInterceptor(s.l.Desugar(), grpc_zap.WithLevels(customCodeToLevel)),
//			)),
//		)
//
//		if err := s.config.PrometheusRegistry.Register(grpcMetrics); err != nil {
//			return err
//		}
//
//		if err := s.config.PrometheusRegistry.Register(prometheus.NewGaugeFunc(
//			prometheus.GaugeOpts{
//				Name: "discovery_replication_factor",
//				Help: "A metric with with constant value showing the configured replication factor.",
//			},
//			func() float64 { return float64(s.config.NumReplicas) },
//		)); err != nil {
//			return err
//		}
//
//		r, err := registry.New(s.backend, s.config.PrometheusRegistry, s.l, s.config.NumReplicas)
//		if err != nil {
//			return err
//		}
//
//		go r.StartCacheUpdater(ctx, cacheSyncInterval)
//		go r.StartServiceCounterUpdater(ctx, serviceCounterUpdateInterval)
//
//		ns, err := r.ListNamespaces()
//		if err != nil {
//			return err
//		}
//
//		dflt := discovery.DefaultNamespace()
//
//		if ns.Index(dflt.Name) < 0 {
//			s.l.Infow("creating default namespace", "name", dflt.Name)
//
//			if _, err := r.RegisterNamespace(*dflt); err != nil {
//				return err
//			}
//		}
//
//		a := &API{
//			r:            r,
//			tokenHandler: tokenHandler,
//		}
//
//		// discoveryv1.RegisterServerAPIServer(s.grpcServer, a)
//		// discoveryv1.RegisterServiceAPIServer(s.grpcServer, a)
//		// discoveryv1.RegisterNamespaceAPIServer(s.grpcServer, a)
//		// discoveryv1.RegisterTokenAPIServer(s.grpcServer, a)
//
//		// grpc reflection support
//		reflection.Register(s.grpcServer)
//
//		listener, err := net.Listen("tcp", s.config.GRPCListenAddr)
//		if err != nil {
//			return err
//		}
//
//		return s.grpcServer.Serve(listener)
//	}

func (s *Server) stop() {
	s.ready.Swap(false)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.l.Error("shutdown connect server", slog.Any("err", err))
	}

	s.l.Info("http server stopped")
}

func (s *Server) readyHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		if s.ready.Load() {
			w.WriteHeader(http.StatusNoContent)

			return
		}

		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
