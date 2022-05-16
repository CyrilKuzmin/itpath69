package web

import (
	"net/http"
	"time"

	"github.com/CyrilKuzmin/itpath69/internal/domain/module"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func (w *Web) getUsernameIfAny(c echo.Context) string {
	sess, err := w.session.Get(c.Request(), "session")
	if err != nil {
		return ""
	}
	if sess.Values["username"] != nil {
		return sess.Values["username"].(string)
	}
	return ""
}

func (w *Web) setUserSession(c echo.Context, username string) {
	sess, _ := w.session.Get(c.Request(), "session")
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["username"] = username
	sess.Save(c.Request(), c.Response())
}

func (w *Web) deleteUserSession(c echo.Context, sess *sessions.Session) {
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: -1,
	}
	sess.Save(c.Request(), c.Response())
	// Ensure that it will work everywhere
	c.SetCookie(&http.Cookie{
		Name:    "session",
		Value:   "",
		MaxAge:  -1,
		Path:    "/",
		Domain:  "",
		Expires: time.Now().Add(-24 * time.Hour),
	})
}

func shiftMetas(in []module.ModuleDTO) {
	for in[len(in)-ModulesPerRow].Id != 1 {
		for k := 0; k < ModulesPerRow; k++ {
			less := in[len(in)-1]
			for i := len(in) - 1; i >= 0; i-- {
				if i == 0 {
					in[i] = less
					continue
				}
				in[i] = in[i-1]
			}
		}
	}
}
