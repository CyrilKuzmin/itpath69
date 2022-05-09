package web

import (
	"context"
	"net/http"
	"os"

	"github.com/CyrilKuzmin/itpath69/internal/domain/comment"
	"github.com/CyrilKuzmin/itpath69/internal/domain/module"
	"github.com/CyrilKuzmin/itpath69/internal/domain/tests"
	"github.com/CyrilKuzmin/itpath69/internal/domain/users"
	"github.com/brpaz/echozap"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

type Web struct {
	session        sessions.Store
	log            *zap.Logger
	userService    users.Service
	moduleService  module.Service
	commentService comment.Service
	testsService   tests.Service
	e              *echo.Echo
}

func NewWeb(log *zap.Logger,
	sessionStore sessions.Store,
	us users.Service,
	ms module.Service,
	cs comment.Service,
	ts tests.Service) *Web {
	w := &Web{
		session:        sessionStore,
		log:            log,
		userService:    us,
		moduleService:  ms,
		commentService: cs,
		testsService:   ts,
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

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	he, ok := err.(*echo.HTTPError)
	if !ok {
		c.JSON(code, err.Error())
	}
	code = he.Code
	// we'll handle 400 and 401 by the client and display an error
	if code < 403 {
		c.Logger().Error(err)
		c.String(code, err.Error())
		return
	}
	c.Logger().Error(err)
	if err := c.Render(code, "error.html", map[string]interface{}{
		"Code":  he.Code,
		"Error": he.Message,
	}); err != nil {
		c.Logger().Error(err)
	}
}

func (w *Web) initHandlers() {
	w.e.Use(session.Middleware(w.session))
	w.e.HTTPErrorHandler = customHTTPErrorHandler

	// render pages
	w.e.GET("/", w.indexHandler)          // render main page
	w.e.GET("/login", w.loginPageHandler) // render login form page
	w.e.GET("/lk", w.lkHandler)           // restricted render LK page with modules previews
	w.e.GET("/module", w.moduleHandler)   // restricted render module page with comments
	w.e.GET("/testing", w.testingHandler) // restricted render testing page

	// User API
	w.e.POST("/user/login", w.loginHandler)
	w.e.POST("/user/logout", w.logoutHandler)
	w.e.POST("/user/register", w.registerHandler)
	w.e.GET("/user/tests", w.listTestsHandler)

	// Comments API
	w.e.GET("/comment", w.getCommentsHandler)
	w.e.POST("/comment", w.newCommentHandler)
	w.e.PUT("/comment", w.updateCommentHandler)
	w.e.DELETE("/comment", w.deleteCommentHandler)

	// Tests API
	w.e.GET("/test", w.getTestHandler)
	w.e.POST("/test", w.checkTestHandler)

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
