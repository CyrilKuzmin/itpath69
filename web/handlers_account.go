package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (w *Web) accountHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	// get user
	user, err := w.srv.GetUserByName(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
	}
	tests, err := w.srv.ListTestsByUser(c.Request().Context(), user.Id)
	if err != nil {
		return errInternal(err)
	}
	return c.Render(http.StatusOK, "account.html", map[string]interface{}{
		"Username": username,
		"User":     user,
		"Tests":    tests,
	})
}
