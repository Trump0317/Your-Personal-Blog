package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ypb/your-personal-blog/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultAddress         = ":8080"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

// Server HTTP Server 实例。
type Server struct {
	ctx    context.Context
	eg     *errgroup.Group
	server *http.Server

	Engine *gin.Engine
	notify chan error

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration

	logger logger.Interface
}

// New 创建 Server。
func New(l logger.Interface, opts ...Option) *Server {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1)

	s := &Server{
		ctx:             ctx,
		eg:              group,
		Engine:          gin.New(),
		notify:          make(chan error, 1),
		address:         _defaultAddress,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
		logger:          l,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      s.Engine,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}

	return s
}

// Start 启动 HTTP Server。
func (s *Server) Start() {
	s.eg.Go(func() error {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.notify <- err

			close(s.notify)

			return err
		}

		return nil
	})

	s.logger.Info("http server - Server - Started")
}

// Notify 返回服务启动与运行错误通道。
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown 优雅关闭 HTTP Server。
func (s *Server) Shutdown() error {
	var shutdownErrors []error

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err, "http server - Server - Shutdown - s.server.Shutdown")

		shutdownErrors = append(shutdownErrors, err)
	}

	err = s.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err, "http server - Server - Shutdown - s.eg.Wait")

		shutdownErrors = append(shutdownErrors, err)
	}

	s.logger.Info("http server - Server - Shutdown")

	return errors.Join(shutdownErrors...)
}
