package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/CyrilKuzmin/itpath69/models"
	"github.com/CyrilKuzmin/itpath69/store"
	"github.com/labstack/echo/v4"
)

var startModulesAmount = 4
var modulesMetaInRow = 4

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
	openedModules := len(user.Modules)
	completedModules := countCompletedModules(user.Modules)
	modulesMeta, err := s.st.GetModulesMeta(c.Request().Context(), len(user.Modules))
	if err != nil {
		return errInternal(err)
	}
	for i := 0; i < len(modulesMeta); i++ {
		modulesMeta[i].Completed = !user.Modules[modulesMeta[i].Id].CompletedAt.IsZero()
	}
	rows := previewsToRowsByN(modulesMeta, modulesMetaInRow)
	if len(rows) > 1 {
		// swap rows in desc order
		for i := 0; i < len(rows)/2; i++ {
			rows[i], rows[len(rows)-1-i] = rows[len(rows)-1-i], rows[i]
		}
	}

	return c.Render(http.StatusOK, "lk.html", map[string]interface{}{
		"User":             user,
		"Username":         user.Username, // for navbar
		"Rows":             rows,
		"ModulesTotal":     s.cm.ModulesTotal, // for statistoc
		"ModulesOpened":    openedModules,
		"ModulesCompleted": completedModules,
	})
}

func (s *App) moduleHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := s.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	user, err := s.st.GetUser(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
	}

	idParam := c.QueryParam("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return errInternal(err)
	}
	if id > len(user.Modules) {
		c.String(http.StatusForbidden, "")
	}
	module, err := s.st.GetModule(c.Request().Context(), id)
	if err != nil {
		return errInternal(err)
	}
	data := make([]template.HTML, len(module.Data))
	for i, p := range module.Data {
		data[i] = template.HTML(p.Data)
	}
	return c.Render(http.StatusOK, "module.html", map[string]interface{}{
		"Username": username,
		"User":     user,
		"Module":   module.Meta,
		"Data":     data,
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
	username := c.FormValue("username")
	password := c.FormValue("password")
	user, err := s.st.SaveUser(c.Request().Context(), username, password)
	if err != nil {
		if store.ErrorIs(err, store.AlreadyExistsErr) {
			return errUserAlreadyExists(username)
		} else {
			return errInternal(err)
		}
	}
	currTime := time.Now()
	for i := 1; i <= startModulesAmount; i++ {
		user.Modules[i] = models.ModuleProgress{CreatedAt: currTime}
	}
	err = s.st.UpdateProgress(c.Request().Context(), username, user.Modules)
	if err != nil {
		return errInternal(err)
	}
	s.setUserSession(c, username)
	return c.String(http.StatusOK, "OK")
}

func (s *App) giveMeModules(c echo.Context) error {
	username := s.getUsernameIfAny(c)
	user, _ := s.st.GetUser(c.Request().Context(), username)
	currTime := time.Now()
	for i := len(user.Modules); i <= len(user.Modules)+startModulesAmount; i++ {
		if _, found := user.Modules[i]; found {
			continue
		}
		user.Modules[i] = models.ModuleProgress{CreatedAt: currTime}
	}
	return s.st.UpdateProgress(c.Request().Context(), username, user.Modules)
}

func (s *App) completeModule(c echo.Context) error {
	username := s.getUsernameIfAny(c)
	idParam := c.QueryParam("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return errInternal(err)
	}
	currTime := time.Now()
	user, _ := s.st.GetUser(c.Request().Context(), username)
	opened := len(user.Modules)
	created := user.Modules[id].CreatedAt
	user.Modules[id] = models.ModuleProgress{CreatedAt: created, CompletedAt: currTime}
	completedOnStage := 0
	for i := opened - startModulesAmount + 1; i <= opened; i++ {
		if !user.Modules[i].CompletedAt.IsZero() {
			completedOnStage++
		}
	}
	if completedOnStage > 2 {
		for i := len(user.Modules) + 1; i <= opened+startModulesAmount; i++ {
			user.Modules[i] = models.ModuleProgress{CreatedAt: currTime}
		}
	}
	fmt.Println("let's update", user.Modules)
	return s.st.UpdateProgress(c.Request().Context(), username, user.Modules)
}
