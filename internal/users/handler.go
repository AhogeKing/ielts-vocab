package users

import (
	"Ielts-vocab/internal/common"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	s *Service
}

func (h *Handler) Login(c *gin.Context) {
	var loginJson LoginRequest
	if err := c.ShouldBindJSON(&loginJson); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error(), "请求体错误")
		return
	}

	user, err := h.s.Login(loginJson)
	if err != nil {
		var invalidCredentials *ErrInvalidCredentials
		var accountDisabled *ErrAccountDisabled

		switch {
		case errors.As(err, &invalidCredentials):
			common.Fail(c, http.StatusUnauthorized, err.Error(), "用户名或密码错误")
			return
		case errors.As(err, &accountDisabled):
			common.Fail(c, http.StatusForbidden, err.Error(), "账号已被禁用")
			return
		default:
			common.Fail(c, http.StatusInternalServerError, err.Error(), "服务器内部错误")
			return
		}
	}
	_ = user
	common.OK(c, "登录成功")
}

func Register(c *gin.Context) {}
