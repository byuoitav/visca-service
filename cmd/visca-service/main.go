package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/byuoitav/visca"
	"github.com/byuoitav/visca-service/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		port     int
		logLevel int8
	)

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.Int8VarP(&logLevel, "log-level", "L", 0, "level to log at. refer to https://godoc.org/go.uber.org/zap/zapcore#Level for options")
	pflag.Parse()

	// build the logger
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(logLevel)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json", EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "@",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	plain, err := config.Build()
	if err != nil {
		fmt.Printf("unable to build logger: %s", err)
		os.Exit(1)
	}

	sugared := plain.Sugar()
	cameras := &sync.Map{}

	handlers := handlers.Handlers{
		CreateCamera: func(ctx context.Context, addr string) (*visca.Camera, error) {
			if cam, ok := cameras.Load(addr); ok {
				return cam.(*visca.Camera), nil
			}

			cam := visca.New(addr, visca.WithLogger(sugared.Named(addr)))
			cameras.Store(addr, cam)
			return cam, nil
		},
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet},
		MaxAge:       7200, // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
	}))

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	e.GET("/:address/pantilt/up", handlers.TiltUp)
	e.GET("/:address/pantilt/down", handlers.TiltDown)
	e.GET("/:address/pantilt/left", handlers.PanLeft)
	e.GET("/:address/pantilt/right", handlers.PanRight)
	e.GET("/:address/pantilt/stop", handlers.PanTiltStop)

	e.GET("/:address/zoom/in", handlers.ZoomIn)
	e.GET("/:address/zoom/out", handlers.ZoomOut)
	e.GET("/:address/zoom/stop", handlers.ZoomStop)

	e.GET("/:address/memory/recall/:channel", handlers.MemoryRecall)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		sugared.Fatalf("unable to bind listener: %s", err)
	}

	sugared.Infof("Starting server on %s", lis.Addr().String())
	err = e.Server.Serve(lis)
	switch {
	case errors.Is(err, http.ErrServerClosed):
	case err != nil:
		sugared.Fatalf("failed to serve: %s", err)
	}
}
