package user

import "github.com/go-playground/validator/v10"

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Email    string `json:"email" binding:"omitempty,email,max=255"`
}

func registerJSONFieldName(field string) string {
	switch field {
	case "Username":
		return "username"
	case "Password":
		return "password"
	case "Email":
		return "email"
	default:
		return field
	}
}

type (
	LoginRequest struct {
		Username string `json:"username" binding:"required,max=50"`
		Password string `json:"password" binding:"required,max=72"`
	}

	LoginResponse struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
)

func loginJSONFieldName(field string) string {
	switch field {
	case "Username":
		return "username"
	case "Password":
		return "password"
	default:
		return field
	}
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "不能为空"
	case "min":
		return "长度不少于 " + fe.Param() + " 个字符"
	case "max":
		return "长度不超过 " + fe.Param() + " 个字符"
	default:
		return "格式不正确"
	}
}
