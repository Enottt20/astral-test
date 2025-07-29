package handler

import (
	"net/http"
	"strconv"

	"github.com/Enottt20/astral-test/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UploadDocument godoc
// @Summary      Upload a document
// @Description  Uploads a document with metadata
// @Tags         documents
// @Accept       multipart/form-data
// @Produce      json
// @Param        token query string true "Authorization token"
// @Param        file formData file false "File"
// @Param        meta formData string false "Metadata"
// @Param        json formData string false "Extra JSON"
// @Success      200 {object} domain.DataResponse
// @Failure      500 {object} domain.ErrorResponse
// @Router       /api/docs [post]
func (e *Endpoint) UploadDocument(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	isFileLoaded := true
	if err != nil {
		isFileLoaded = false
	}
	defer file.Close()

	meta := c.PostForm("meta")
	jsonData := c.PostForm("json")
	token := c.GetString("token")

	doc, err := e.services.Documents.Upload(c.Request.Context(), token, meta, jsonData, file, header, isFileLoaded)
	if err != nil {
		logrus.Errorf("Failed to upload document (internal error): %s", err.Error())
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error: domain.ErrorInfo{Code: 500, Text: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, domain.DataResponse{
		Data: doc,
	})
}

// GetDocuments godoc
// @Summary      Get documents
// @Description  Retrieves list of documents
// @Tags         documents
// @Produce      json
// @Param        token query string true "Authorization token"
// @Param        login query string false "Filter by login"
// @Param        key query string false "Metadata key"
// @Param        value query string false "Metadata value"
// @Param        limit query int false "Limit number of results"
// @Success      200 {object} domain.DataResponse
// @Failure      500 {object} domain.ErrorResponse
// @Router       /api/docs [get]
func (e *Endpoint) GetDocuments(c *gin.Context) {
	token := c.GetString("token")
	login := c.Query("login")
	key := c.Query("key")
	value := c.Query("value")
	limit, _ := strconv.Atoi(c.Query("limit"))

	docs, err := e.services.Documents.GetAll(c.Request.Context(), token, login, key, value, limit)
	if err != nil {
		logrus.Errorf("Failed to get documents (internal error): %s", err.Error())
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error: domain.ErrorInfo{Code: 500, Text: err.Error()},
		})
		return
	}
	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
		return
	}
	c.JSON(http.StatusOK, domain.DataResponse{
		Data: gin.H{"docs": docs},
	})
}

// GetDocument godoc
// @Summary      Get document by ID
// @Description  Returns a document or its file
// @Tags         documents
// @Produce      json
// @Param        token query string true "Authorization token"
// @Param        id path string true "Document ID"
// @Success      200 {object} domain.DataResponse
// @Failure      500 {object} domain.ErrorResponse
// @Router       /api/docs/{id} [get]
func (e *Endpoint) GetDocument(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")

	doc, fileData, err := e.services.Documents.GetByID(c.Request.Context(), token, id)
	if err != nil {
		logrus.Errorf("Failed to get document (internal error): %s", err.Error())
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error: domain.ErrorInfo{Code: 500, Text: err.Error()},
		})
		return
	}
	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
		return
	}
	if doc.File {
		c.Data(http.StatusOK, doc.Mime, fileData)
	} else {
		c.JSON(http.StatusOK, domain.DataResponse{
			Data: string(fileData),
		})
	}
}

// DeleteDocument godoc
// @Summary      Delete a document
// @Description  Deletes document by ID
// @Tags         documents
// @Produce      json
// @Param        token query string true "Authorization token"
// @Param        id path string true "Document ID"
// @Success      200 {object} domain.Response
// @Failure      500 {object} domain.ErrorResponse
// @Router       /api/docs/{id} [delete]
func (e *Endpoint) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	token := c.GetString("token")

	err := e.services.Documents.Delete(c.Request.Context(), token, id)
	if err != nil {
		logrus.Errorf("Failed to delete document (internal error): %s", err.Error())
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error: domain.ErrorInfo{Code: 500, Text: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, domain.Response{
		Response: gin.H{id: true},
	})
}
