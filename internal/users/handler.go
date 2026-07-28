package users

import (
	"Ielts-vocab/internal/common"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	s *Service
}

func (h *Handler) Login(c *gin.Context) {
	var loginJSON LoginRequest
	if err := c.ShouldBindJSON(&loginJSON); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			fields := make(map[string]string, len(validationErrs))

			for _, fe := range validationErrs {
				field := loginJSONFieldName(fe.Field())
				fields[field] = validationMessage(fe)
			}

			common.FailWithFields(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不合法", fields)
		} else {
			common.Fail(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", "请求体错误")
		}
		return
	}

	user, err := h.s.Login(loginJSON)
	if err != nil {
		var invalidCredentials *ErrInvalidCredentials
		var accountDisabled *ErrAccountDisabled

		switch {
		case errors.As(err, &invalidCredentials):
			common.Fail(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "用户名或密码错误")
			return
		case errors.As(err, &accountDisabled):
			common.Fail(c, http.StatusForbidden, "ACCOUNT_ALREADY_DISABLED", "账号已被禁用")
			return
		default:
			common.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
			return
		}
	}
	_ = user
	common.OK(c, "登录成功")
}

func (h *Handler) Register(c *gin.Context) {
	var regJSON RegisterRequest
	if err := c.ShouldBindJSON(&regJSON); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			fields := make(map[string]string, len(validationErrs))

			for _, fe := range validationErrs {
				field := registerJSONName(fe.Field())
				fields[field] = validationMessage(fe)
			}

			common.FailWithFields(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不合法", fields)
		} else {
			common.Fail(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", "请求体错误")
		}
		return
	}

	if err := h.s.Register(regJSON); err != nil {
		var usernameAlreadyExist *ErrUsernameAlreadyExist
		if errors.As(err, &usernameAlreadyExist) {
			common.Fail(c, http.StatusConflict, "USERNAME_ALREADY_EXISTS", "用户名已存在")
			return
		}

		common.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
		return
	}

	common.Created(c, "注册成功")
}
