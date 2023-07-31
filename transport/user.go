package transport

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ynuraddi/t-tsarka/model"
)

type userCreateRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type userCreateResponce struct {
	ID int64 `json:"id"`
}

func (s *Server) userCreate(c echo.Context) error {
	var req userCreateRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error("failed bind request", err)
		return c.JSON(http.StatusUnprocessableEntity, errorResponce{err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		s.logger.Error("bad request", err)
		return c.JSON(http.StatusBadRequest, errorResponce{err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := s.service.User.Create(ctx, model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		s.logger.Error("failed create user", err)
		return c.JSON(http.StatusInternalServerError, errorResponce{err.Error()})
	}

	s.logger.Debug("user seccessfully created")
	return c.JSON(http.StatusCreated, userCreateResponce{id})
}

func (s *Server) userGet(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		s.logger.Error("failed convert param", err)
		return c.JSON(http.StatusBadRequest, errorResponce{err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbuser, err := s.service.User.Get(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		s.logger.Info("user is not found")
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		s.logger.Error("failed get user", err)
		return c.JSON(http.StatusInternalServerError, errorResponce{err.Error()})
	}

	return c.JSON(http.StatusOK, dbuser)
}

type userUpdateRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

func (s *Server) userUpdate(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponce{err.Error()})
	}

	var req userUpdateRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error("failed bind request", err)
		return c.JSON(http.StatusUnprocessableEntity, errorResponce{err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		s.logger.Error("bad request", err)
		return c.JSON(http.StatusBadRequest, errorResponce{err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbuser, err := s.service.User.Update(ctx, model.User{
		ID:        int64(id),
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if errors.Is(err, sql.ErrNoRows) {
		s.logger.Info("user is not found")
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		s.logger.Error("failed update user", err)
		return c.JSON(http.StatusInternalServerError, errorResponce{err.Error()})
	}

	return c.JSON(http.StatusOK, dbuser)
}

func (s *Server) userDelete(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponce{err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.service.User.Delete(ctx, int64(id)); errors.Is(err, sql.ErrNoRows) {
		s.logger.Info("user is not found")
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		s.logger.Error("failed delete user", err)
		return c.JSON(http.StatusInternalServerError, errorResponce{err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
