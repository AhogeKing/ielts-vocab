package user

import (
	"errors"
	"ielts-vocab/internal/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req LoginRequest
	)
	if !bindJSON(c, &req, loginJSONFieldName) {
		return
	}

	_, err := h.service.Login(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			common.Fail(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "用户名或密码错误")
			return

		case errors.Is(err, ErrAccountDisabled):
			common.Fail(c, http.StatusForbidden, "ACCOUNT_ALREADY_DISABLED", "账号已被禁用")
			return

		default:
			common.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
			return
		}
	}
	common.OK(c, "登录成功")
}

func (h *Handler) Register(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req RegisterRequest
	)

	if !bindJSON(c, &req, registerJSONFieldName) {
		return
	}

	if err := h.service.Register(ctx, req); err != nil {
		if errors.Is(err, ErrUsernameAlreadyExist) {
			common.Fail(c, http.StatusConflict, "USERNAME_ALREADY_EXISTS", "用户名已存在")
			return
		}

		common.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
		return
	}

	common.Created(c, "注册成功")
}

func bindJSON[T any](c *gin.Context, req *T, jsonFieldName func(string) string) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		if validationErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			fields := make(map[string]string, len(validationErrs))

			for _, fe := range validationErrs {
				fields[jsonFieldName(fe.Field())] = validationMessage(fe)
			}

			common.FailWithFields(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数不合法", fields)
		} else {
			common.Fail(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", "请求体错误")
		}
		return false
	}
	return true
}
