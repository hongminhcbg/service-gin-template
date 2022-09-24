package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hongminhcbg/service-gin-template/src/service"
)

func InitGin(e *gin.Engine, s *service.Service) {
	e.POST("/api/v1/users", s.CreateUser)
}
