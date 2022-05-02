package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	mongosessions "github.com/2-72/gorilla-sessions-mongodb"

	"github.com/CyrilKuzmin/itpath69/config"
	"github.com/CyrilKuzmin/itpath69/content"
	"github.com/CyrilKuzmin/itpath69/internal/domain/module"
	"github.com/CyrilKuzmin/itpath69/internal/domain/users"
	"github.com/CyrilKuzmin/itpath69/store/mongostorage"
	"github.com/CyrilKuzmin/itpath69/web"

	"go.uber.org/zap"
)

type App struct {
	config         *config.Config
	web            *web.Web
	log            *zap.Logger
	contentManager *content.ContentManager
}

func NewApp(conf *config.Config) *App {
	// init logger first
	log := initLogger(conf.Env)
	// init store
	st := mongostorage.NewMongo(log, conf.Mongo.URI, conf.Mongo.Database)
	// init sessions storage
	session, err := mongosessions.NewMongoDBStore(st.Sessions, []byte(conf.Secret))
	if err != nil {
		log.Error("cannot init session storage", zap.Error(err))
	}
	// init echo server with its dependencies
	us := users.NewService(st, log)
	ms := module.NewService(st, log)
	webserver := web.NewWeb(log, session, us, ms)
	// init content manager
	cm := content.NewContentManager(ms, log)
	// create App and init handlers/middlewares
	s := &App{
		config:         conf,
		web:            webserver,
		log:            log,
		contentManager: cm,
	}
	cm.UpdateContentinStorage()
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

// Start server
func (s *App) Start() {
	addr := fmt.Sprintf("%v:%v", s.config.Host, s.config.Port)
	go func() {
		if err := s.web.Start(addr); err != nil && err != http.ErrServerClosed {
			s.log.Fatal("shutting down the server")
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.web.Shutdown(ctx); err != nil {
		s.log.Fatal("cannot shutdown properly", zap.Error(err))
	}
	//s.Close(ctx)
}
