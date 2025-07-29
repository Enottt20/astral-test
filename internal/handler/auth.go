package handler

import (
	"net/http"

	"github.com/Enottt20/astral-test/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Register godoc
// @Summary      Register a new user
// @Description  Registers user and returns login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body domain.RegisterRequest true "User credentials"
// @Success      200 {object} domain.Response
// @Failure      400 {object} domain.ErrorResponse
// @Failure      500 {object} domain.ErrorResponse
// @Router       /api/register [post]
func (e *Endpoint) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Failed to register user (invalid request): %s", err.Error())
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: domain.ErrorInfo{Code: 400, Text: "Invalid request format"},
		})
		return
	}
	login, err := e.services.Auth.Register(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to register user (internal error): %s", err.Error())
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error: domain.ErrorInfo{Code: 500, Text: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, domain.Response{
		Response: gin.H{"login": login},
	})
}

// Authenticate godoc
// @Summary      Authenticate user
// @Description  Authenticates user and returns JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body domain.AuthRequest true "User credentials"
// @Success      200 {object} domain.Response
// @Failure      400 {object} domain.ErrorResponse
// @Failure      401 {object} domain.ErrorResponse
// @Router       /api/auth [post]
func (e *Endpoint) Authenticate(c *gin.Context) {
	var req domain.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("Failed to authenticate user (invalid request): %s", err.Error())
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: domain.ErrorInfo{Code: 400, Text: "Invalid request format"},
		})
		return
	}
	token, err := e.services.Auth.Authenticate(c.Request.Context(), req)
	if err != nil {
		logrus.Errorf("Failed to authenticate user (internal error): %s", err.Error())
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Error: domain.ErrorInfo{Code: 401, Text: "Invalid credentials"},
		})
		return
	}
	c.JSON(http.StatusOK, domain.Response{
		Response: gin.H{"token": token},
	})
}

func (e *Endpoint) Logout(c *gin.Context) {
	token := c.Param("token")
	if err := e.services.Auth.Logout(c.Request.Context(), token); err != nil {
		logrus.Errorf("Failed to logout user (internal error): %s", err.Error())
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error: domain.ErrorInfo{Code: 500, Text: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, domain.Response{
		Response: gin.H{token: true},
	})
}
