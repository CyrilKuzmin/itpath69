package web

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

const MaximumCommentLength = 10000

func (w *Web) getCommentsHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	moduleIdValue := c.QueryParam("module_id")
	moduleId, err := strconv.Atoi(moduleIdValue)
	if err != nil {
		return errInternal(err)
	}
	comments, err := w.srv.ListCommentsByModule(c.Request().Context(), username, moduleId)
	if err != nil {
		return errInternal(err)
	}
	return c.JSON(http.StatusOK, comments)
}

func (w *Web) newCommentHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	text := c.FormValue("text")
	if len(text) > MaximumCommentLength {
		return errCommentTooLong()
	}
	moduleIdValue := c.FormValue("module_id")
	moduleId, err := strconv.Atoi(moduleIdValue)
	if err != nil {
		return errInternal(err)
	}
	partIdValue := c.FormValue("part_id")
	partId, err := strconv.Atoi(partIdValue)
	if err != nil {
		return errInternal(err)
	}
	// create comment
	comment, err := w.srv.CreateComment(c.Request().Context(), username, text, moduleId, partId)
	if err != nil {
		return errInternal(err)
	}
	// send OK
	return c.JSON(http.StatusOK, comment)
}

func (w *Web) updateCommentHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	commentId := c.FormValue("id")
	text := c.FormValue("text")
	if len(text) > MaximumCommentLength {
		return errCommentTooLong()
	}
	// update comment
	newComment, err := w.srv.UpdateComment(c.Request().Context(), username, commentId, text)
	if err != nil {
		return errInternal(err)
	}
	// send OK
	return c.JSON(http.StatusOK, newComment)
}

func (w *Web) deleteCommentHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	commentId := c.QueryParam("id")
	if commentId == "" {
		return errBadRequest()
	}
	// delete comment
	err := w.srv.DeleteCommentByID(c.Request().Context(), username, commentId)
	if err != nil {
		return errInternal(err)
	}
	// send OK
	return c.String(http.StatusOK, "OK")
}
