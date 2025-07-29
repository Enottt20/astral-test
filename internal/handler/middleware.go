package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Enottt20/astral-test/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BodyWithToken struct {
	Token string `json:"token" binding:"required"`
}

func (e *Endpoint) Middleware(c *gin.Context) {
	var token string
	var req BodyWithToken
	var metaData domain.DocumentMeta
	meta := c.PostForm("meta")
	if err := json.Unmarshal([]byte(meta), &metaData); err == nil {
		token = metaData.Token
	}
	if err := c.ShouldBindJSON(&req); err == nil {
		token = req.Token
	}
	if c.Query("token") != "" {
		token = c.Query("token")
	}
	if token == "" {
		logrus.Error("Middleware error: token required")
		c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorResponse{
			Error: domain.ErrorInfo{
				Text: "Authorization token required",
				Code: 401,
			},
		})
		return
	}
	valid, err := e.services.Auth.ValidateToken(c.Request.Context(), token)
	if err != nil || !valid {
		logrus.Error("Middleware error: token is invalid")
		c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorResponse{
			Error: domain.ErrorInfo{
				Text: "Invalid token",
				Code: 401,
			},
		})
		return
	}
	c.Set("token", token)
	c.Next()
}
