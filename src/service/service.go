package service

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hongminhcbg/service-gin-template/config"
	"github.com/hongminhcbg/service-gin-template/src/models"
	"github.com/hongminhcbg/service-gin-template/src/store"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type Service struct {
	store *store.UserStore
}

func NewService(cfg *config.Config, store *store.UserStore, r *redis.Client) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) createNewUser(ctx *gin.Context, req *models.CreateUserRequest) {
	if req.ReqId == "" {
		req.ReqId = fmt.Sprint(time.Now().UnixMilli())
	}

	r := models.User{
		Name:      req.Name,
		ReqId:     req.ReqId,
		RetryTime: 0,
		Status:    models.USER_INIT,
	}

	err := s.store.Save(ctx.Request.Context(), &r)
	if err != nil {
		log.Println(err, "save to record error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	log.Println("create new user success")
	ctx.JSON(http.StatusOK, r)
}

func (s *Service) handleExistedReqId(ctx *gin.Context, u *models.User, req *models.CreateUserRequest) {
	ctx.JSON(http.StatusOK, u)
}

func (s *Service) CreateUser(ctx *gin.Context) {
	var req models.CreateUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if req.ReqId == "" {
		s.createNewUser(ctx, &req)
		return
	}

	u, _ := s.store.GetByReqId(ctx.Request.Context(), req.ReqId)
	if u == nil {
		s.createNewUser(ctx, &req)
		return
	}

	s.handleExistedReqId(ctx, u, &req)
}