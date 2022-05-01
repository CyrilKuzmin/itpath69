package server

import (
	"net/http"
	"time"

	"github.com/CyrilKuzmin/itpath69/models"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func (s *App) getUsernameIfAny(c echo.Context) string {
	sess, err := s.session.Get(c.Request(), "session")
	if err != nil {
		return ""
	}
	if sess.Values["username"] != nil {
		return sess.Values["username"].(string)
	}
	return ""
}

func (s *App) setUserSession(c echo.Context, username string) {
	sess, _ := s.session.Get(c.Request(), "session")
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["username"] = username
	sess.Save(c.Request(), c.Response())
}

func (s *App) deleteUserSession(c echo.Context, sess *sessions.Session) {
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

func previewsToRowsByN(modulesPreviews []models.ModuleMeta, n int) []ModulesRow {
	rows := make([]ModulesRow, 0)
	if len(modulesPreviews) <= n {
		rows = append(rows, ModulesRow{Modules: modulesPreviews})
		return rows
	}
	rows = append(rows, ModulesRow{Modules: modulesPreviews[:n]})
	rows = append(rows, previewsToRowsByN(modulesPreviews[n:], n)...)
	return rows
}
