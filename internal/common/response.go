package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code   string            `json:"code"`
	Msg    string            `json:"msg,omitempty"`
	Fields map[string]string `json:"fields,omitempty"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// OK sends a success response
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created sends a successful resource-creation response.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// Fail sends an error response
func Fail(c *gin.Context, status int, code, msg string) {
	c.JSON(status, Response{
		Success: false,
		Error:   &ErrorInfo{Code: code, Msg: msg},
	})
}

func FailWithFields(c *gin.Context, status int, code, msg string, fields map[string]string) {
	c.JSON(status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:   code,
			Msg:    msg,
			Fields: fields,
		},
	})
}
