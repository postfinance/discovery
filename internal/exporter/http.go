package exporter

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	httpStopTimeout   = 10 * time.Second
	readHeaderTimeout = 10 * time.Second
)

func (e *Exporter) startHTTP() error {
	e.log.Infow("starting http server")

	mux := http.NewServeMux()

	r, ok := e.config.PrometheusRegistry.(prometheus.Gatherer)
	if !ok {
		panic("interface is not prometheus.Registry")
	}

	if err := e.config.PrometheusRegistry.Register(e.serviceWatchEvents); err != nil {
		return err
	}

	mux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	e.httpServer = &http.Server{
		Addr:        e.config.HTTPListenAddr,
		Handler:     mux,
		ReadTimeout: readHeaderTimeout,
	}

	if err := e.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (e *Exporter) stopHTTP() error {
	defer e.log.Info("http server stopped")

	if e.httpServer == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), httpStopTimeout)
	defer cancel()

	return e.httpServer.Shutdown(ctx)
}
