package server

import (
	"net/http"

	"github.com/CyrilKuzmin/itpath69/models"
	"github.com/CyrilKuzmin/itpath69/store"
	"github.com/labstack/echo/v4"
)

var startModulesAmount = 4
var modulesPreviewsInRow = 4

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

type ModulesRow struct {
	Modules []models.ModuleMeta
}

func (s *App) lkHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := s.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	// get the list of opened modules and show them
	user, err := s.st.GetUser(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
	}
	modulesPreviews, err := s.st.GetModulesMeta(c.Request().Context(), user.ModulesOpened)
	if err != nil {
		return errInternal(err)
	}
	rows := previewsToRowsByN(modulesPreviews, modulesPreviewsInRow)
	if len(rows) > 1 {
		// swap rows in desc order
		for i := 0; i < len(rows)/2; i++ {
			rows[i], rows[len(rows)-1-i] = rows[len(rows)-1-i], rows[i]
		}
	}

	return c.Render(http.StatusOK, "lk.html", map[string]interface{}{
		"Username": username,
		"Rows":     rows,
	})
}

func (s *App) moduleHandler(c echo.Context) error {
	sess, err := s.session.Get(c.Request(), "session")
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	return c.Render(http.StatusOK, "module.html", map[string]interface{}{
		"Username": sess.Values["username"],
		"Id":       c.QueryParam("id"),
	})
}

func (s *App) loginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	_, err := s.st.CheckUserPassword(c.Request().Context(), username, password)
	if err != nil {
		return errLoginFailed()
	}
	s.setUserSession(c, username)
	return c.String(http.StatusOK, "OK")
}

func (s *App) logoutHandler(c echo.Context) error {
	sess, err := s.session.Get(c.Request(), "session")
	if err != nil || sess.ID == "" {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
	s.deleteUserSession(c, sess)
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func (s *App) registerHandler(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("password")
	err := s.st.SaveUser(c.Request().Context(), username, password)
	if err != nil {
		if store.ErrorIs(err, store.AlreadyExistsErr) {
			return errUserAlreadyExists(username)
		} else {
			return errInternal(err)
		}
	}
	err = s.st.OpenModules(c.Request().Context(), username, startModulesAmount)
	if err != nil {
		return errInternal(err)
	}
	s.setUserSession(c, username)
	return c.String(http.StatusOK, "OK")
}
