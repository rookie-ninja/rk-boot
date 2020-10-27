package rk_gin

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/rookie-ninja/rk-boot/sw"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
	"path"
	"strconv"
	"syscall"
)

type GinServerEntry struct {
	logger       *zap.Logger
	router       *gin.Engine
	server       *http.Server
	name         string
	port         uint64
	interceptors []gin.HandlerFunc
	sw           *rk_sw.SWEntry
}

type GinEntryOption func(*GinServerEntry)

func WithRouter(router *gin.Engine) GinEntryOption {
	return func(entry *GinServerEntry) {
		entry.router = router
	}
}

func WithSWEntry(sw *rk_sw.SWEntry) GinEntryOption {
	return func(entry *GinServerEntry) {
		entry.sw = sw
	}
}

func WithPort(port uint64) GinEntryOption {
	return func(entry *GinServerEntry) {
		entry.port = port
	}
}

func WithName(name string) GinEntryOption {
	return func(entry *GinServerEntry) {
		entry.name = name
	}
}

func NewGinServerEntry(opts ...GinEntryOption) *GinServerEntry {
	entry := &GinServerEntry{
		logger: zap.NewNop(),
	}

	for i := range opts {
		opts[i](entry)
	}

	if len(entry.name) < 1 {
		entry.name = "gin-server-" + strconv.FormatUint(entry.port, 10)
	}

	if entry.interceptors == nil {
		entry.interceptors = make([]gin.HandlerFunc, 0)
	}

	if entry.router == nil {
		gin.SetMode(gin.ReleaseMode)
		entry.router = gin.New()
	}

	if entry.sw != nil {
		entry.router.GET(path.Join(entry.sw.GetPath(), "*any"), entry.sw.GinHandler())
		entry.router.GET("/swagger/*any", entry.sw.GinFileHandler())
	}

	if len(entry.interceptors) > 0 {
		entry.router.Use(entry.interceptors...)
	}

	endpoint := "0.0.0.0:" + strconv.FormatUint(entry.port, 10)

	entry.server = &http.Server{
		Addr:    endpoint,
		Handler: entry.router,
	}

	return entry
}

func (entry *GinServerEntry) GetName() string {
	return entry.name
}

func (entry *GinServerEntry) GetPort() uint64 {
	return entry.port
}

func (entry *GinServerEntry) GetSWEntry() *rk_sw.SWEntry {
	return entry.sw
}

func (entry *GinServerEntry) GetServer() *http.Server {
	return entry.server
}

func (entry *GinServerEntry) GetRouter() *gin.Engine {
	return entry.router
}

func (entry *GinServerEntry) AddInterceptor(inters ...gin.HandlerFunc) {
	entry.interceptors = append(entry.interceptors, inters...)
}

func (entry *GinServerEntry) Start(logger *zap.Logger) {
	if logger == nil {
		logger = zap.NewNop()
	}

	go func(entry *GinServerEntry) {
		logger.Info("starting gin server",
			zap.Uint64("gin_port", entry.port),
			zap.String("name", entry.name))
		if entry.sw != nil {
			logger.Info("starting swagger",
				zap.Uint64("sw_port", entry.sw.GetSWPort()),
				zap.String("sw_path", entry.sw.GetPath()))
		}
		if err := entry.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("err while serving gin listener",
				zap.Uint64("gin_port", entry.port),
				zap.String("name", entry.name),
				zap.Error(err))
			shutdownWithError(err)
		}
	}(entry)
}

func (entry *GinServerEntry) Stop(logger *zap.Logger) {
	if entry.router != nil {
		if logger == nil {
			logger = zap.NewNop()
		}

		logger.Info("stopping gin server",
			zap.Uint64("gin_port", entry.port),
			zap.String("name", entry.name))
		if err := entry.server.Shutdown(context.Background()); err != nil {
			logger.Warn("error occurs while stopping gin server",
				zap.Uint64("gin_port", entry.port),
				zap.String("name", entry.name),
				zap.Error(err))
		}
	}
}

func shutdownWithError(err error) {
	glog.Error(err)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
