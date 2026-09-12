// Package response is the single place that shapes every JSON response the
// API sends. Handlers must always go through Success/Created/Paginated/Error
// instead of calling c.JSON directly, so the response shape stays consistent
// across every endpoint.
package response

import "github.com/gin-gonic/gin"

// Envelope is the standard shape of every API response.
type Envelope struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
	Meta    *Meta      `json:"meta,omitempty"`
}

// ErrorBody carries a machine-readable code plus a human-readable message.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta carries response metadata, such as the total count for a list.
type Meta struct {
	Total int `json:"total"`
}

// Success writes a 200 response with the given data.
func Success(c *gin.Context, data any) {
	c.JSON(200, Envelope{Success: true, Data: data})
}

// Created writes a 201 response with the given data.
func Created(c *gin.Context, data any) {
	c.JSON(201, Envelope{Success: true, Data: data})
}

// Paginated writes a 200 response with data plus a total count in Meta.
func Paginated(c *gin.Context, data any, total int) {
	c.JSON(200, Envelope{Success: true, Data: data, Meta: &Meta{Total: total}})
}

// Error writes an error response at the given HTTP status and aborts the
// request, so no further handler/middleware writes to the response.
func Error(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, Envelope{
		Success: false,
		Error:   &ErrorBody{Code: code, Message: message},
	})
}
