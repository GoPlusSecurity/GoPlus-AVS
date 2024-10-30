// AVS 的 HTTP Server 对象
package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"goplus/avs/config"
	"goplus/avs/metrics"
	"goplus/avs/secwaremanager"
	"goplus/avs/state"
	"net/http"
	"time"
)

type Server struct {
	config         config.Config
	logger         logging.Logger
	metricsIntf    metrics.AvsMetricsInterface
	secwareManager *secwaremanager.SecwareManager
	stateIntf      state.AvsStateInterface

	secwareAccessorIntf secwaremanager.SecwareAccessorInterface
}

func New(cfg config.Config, metricsIntf metrics.AvsMetricsInterface, secwareManager *secwaremanager.SecwareManager, state state.AvsStateInterface) (*Server, error) {
	return &Server{
		config:         cfg,
		logger:         cfg.Logger,
		metricsIntf:    metricsIntf,
		secwareManager: secwareManager,
		stateIntf:      state,

		secwareAccessorIntf: &secwaremanager.SecwareAccessorImpl{},
	}, nil
}

func (a *Server) Init() error {
	return nil
}

type ginLogger struct {
	logger logging.Logger
}

func (g *ginLogger) Write(p []byte) (n int, err error) {
	g.logger.Info(string(p))
	return len(p), nil
}

func (a *Server) Start(ctx context.Context) error {
	a.logger.Info("AVS Server start")

	gin.SetMode(gin.ReleaseMode)
	route := gin.New()

	gl := &ginLogger{logger: a.logger}
	route.Use(gin.LoggerWithWriter(gl))
	route.Use(gin.Recovery())

	route.POST("/secware/task", a.handleTask)
	route.GET("/avs/ping", a.ping)
	route.GET("/metrics", gin.WrapH(promhttp.Handler()))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.APIPort),
		Handler: route.Handler(),
	}

	apiDone := make(chan struct{})
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error(err.Error())
			close(apiDone)
		}
	}()

	select {
	case <-ctx.Done():
	case <-apiDone:
		a.logger.Info("Server shutdown...")
		return errors.New("server exit unexpectedly")
	}

	shutdownCtx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		a.logger.Error(err.Error())
	}

	<-shutdownCtx.Done()
	a.logger.Info("Server shutdown complete")
	return nil
}
