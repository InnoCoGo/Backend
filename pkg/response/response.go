package response

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

// type statusOkResponse struct {
// 	Status string `json:"status"`
// }

type DataResponse struct {
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, ErrorResponse{message})
}
