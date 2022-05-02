package web

import (
	"context"
	"os"

	"github.com/CyrilKuzmin/itpath69/internal/domain/module"
	"github.com/CyrilKuzmin/itpath69/internal/domain/users"
	"github.com/brpaz/echozap"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

type Web struct {
	session       sessions.Store
	log           *zap.Logger
	userService   users.Service
	moduleService module.Service
	e             *echo.Echo
}

func NewWeb(log *zap.Logger, sessionStore sessions.Store, us users.Service, ms module.Service) *Web {
	w := &Web{
		session:       sessionStore,
		log:           log,
		userService:   us,
		moduleService: ms,
	}
	w.e = initEcho(log)
	w.initHandlers()
	return w
}

func initEcho(l *zap.Logger) *echo.Echo {
	e := echo.New()
	e.Use(echozap.ZapLogger(l))
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetOutput(os.Stdout)
	e.Debug = true
	e.HideBanner = true
	e.Static("/assets", "static/assets")
	e.Renderer = NewTemplateRenderer("static/*.html")
	return e
}

func (w *Web) initHandlers() {
	w.e.Use(session.Middleware(w.session))
	// Unauthenticated routes
	w.e.GET("/", w.indexHandler)
	w.e.GET("/login", w.loginPageHandler)
	w.e.POST("/logout", w.logoutHandler)
	w.e.POST("/register", w.registerHandler)
	w.e.POST("/login", w.loginHandler)
	// restricted
	w.e.GET("/lk", w.lkHandler)
	w.e.GET("/module", w.moduleHandler)
	// temp
	w.e.GET("/more", w.giveMeModules)
	w.e.GET("/complete", w.completeModule)
}

func (w *Web) Start(addr string) error {
	return w.e.Start(addr)
}

func (w *Web) Shutdown(ctx context.Context) error {
	return w.e.Shutdown(ctx)
}
