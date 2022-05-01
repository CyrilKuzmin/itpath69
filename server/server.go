package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	mongosessions "github.com/2-72/gorilla-sessions-mongodb"
	"github.com/CyrilKuzmin/itpath69/config"
	"github.com/CyrilKuzmin/itpath69/store"
	"github.com/CyrilKuzmin/itpath69/store/mongostorage"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

type App struct {
	c       *config.Config
	e       *echo.Echo
	l       *zap.Logger
	st      store.Store
	session *mongosessions.MongoDBStore
}

func NewApp(conf *config.Config) *App {
	// init logger first
	l := initLogger(conf.Env)
	// init echo server
	e := initEcho(l)
	// init store
	st := mongostorage.NewMongo(l, conf.Mongo.URI, conf.Mongo.Database)
	// init sessions storage
	session, err := mongosessions.NewMongoDBStore(st.Sessions, []byte(conf.Secret))
	if err != nil {
		l.Error("cannot init session storage", zap.Error(err))
	}
	// create App and init handlers/middlewares
	s := &App{conf, e, l, st, session}
	s.initHandlers()
	return s
}

func initLogger(env string) *zap.Logger {
	var l *zap.Logger
	if env == "production" {
		l, _ = zap.NewProduction()
	} else {
		l, _ = zap.NewDevelopment()
	}
	return l
}

func initEcho(l *zap.Logger) *echo.Echo {
	e := echo.New()
	e.Use(echozap.ZapLogger(l))
	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetOutput(os.Stdout)
	e.Debug = true
	e.HideBanner = true
	e.Static("/assets", "static/assets")
	e.Renderer = NewTemplateRenderer("static/*.html")
	return e
}

func (s *App) initHandlers() {
	s.e.Use(session.Middleware(s.session))
	// Unauthenticated routes
	s.e.GET("/", s.indexHandler)
	s.e.GET("/login", s.loginPageHandler)
	s.e.POST("/logout", s.logoutHandler)
	s.e.POST("/register", s.registerHandler)
	s.e.POST("/login", s.loginHandler)
	// restricted
	s.e.GET("/lk", s.lkHandler)
	s.e.GET("/module", s.moduleHandler)
}

// Start server
func (s *App) Start() {
	addr := fmt.Sprintf("%v:%v", s.c.Host, s.c.Port)
	go func() {
		if err := s.e.Start(addr); err != nil && err != http.ErrServerClosed {
			s.l.Fatal("shutting down the server")
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.e.Shutdown(ctx); err != nil {
		s.l.Fatal("cannot shutdown properly", zap.Error(err))
	}
	s.st.Close(ctx)
}
