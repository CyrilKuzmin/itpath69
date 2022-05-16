package web

import (
	"net/http"

	"github.com/CyrilKuzmin/itpath69/store"
	"github.com/labstack/echo/v4"
)

func (w *Web) loginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	err := w.srv.CheckUserPassword(c.Request().Context(), username, password)
	if err != nil {
		return errLoginFailed()
	}
	w.setUserSession(c, username)
	return c.String(http.StatusOK, "OK")
}

func (w *Web) logoutHandler(c echo.Context) error {
	sess, err := w.session.Get(c.Request(), "session")
	if err != nil || sess.ID == "" {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
	w.deleteUserSession(c, sess)
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func (w *Web) registerHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	_, err := w.srv.CreateUser(c.Request().Context(), username, password)
	if err != nil {
		if store.ErrorIs(err, store.AlreadyExistsErr) {
			return errUserAlreadyExists(username)
		} else {
			return errInternal(err)
		}
	}
	w.setUserSession(c, username)
	return c.String(http.StatusOK, "OK")
}
