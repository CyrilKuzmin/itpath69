package server

import (
	"net/http"
	"time"

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

func (s *App) indexHandler(c echo.Context) error {
	username := s.getUsernameIfAny(c)
	if username == "" {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
	}
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Username": username,
	})
}

func (s *App) loginPageHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{})
}

func (s *App) lkHandler(c echo.Context) error {
	sess, err := s.session.Get(c.Request(), "session")
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	return c.Render(http.StatusOK, "lk.html", map[string]interface{}{
		"Username": sess.Values["username"],
	})
}

func (s *App) moduleHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "module.html", map[string]interface{}{
		"Username": "Cyrilit",
		"Id":       c.QueryParam("id"),
	})
}

func (s *App) loginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	_, err := s.st.GetUser(c.Request().Context(), username, password)
	if err != nil {
		return echo.ErrUnauthorized
	}
	sess, _ := s.session.Get(c.Request(), "session")
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["username"] = username
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, "/lk")
}

func (s *App) logoutHandler(c echo.Context) error {
	sess, err := s.session.Get(c.Request(), "session")
	if err != nil || sess.ID == "" {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
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
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func (s *App) registerHandler(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("password")
	err := s.st.SaveUser(c.Request().Context(), username, password)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	sess, _ := s.session.Get(c.Request(), "session")
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["username"] = username
	sess.Save(c.Request(), c.Response())
	return c.String(http.StatusOK, "OK")
}
