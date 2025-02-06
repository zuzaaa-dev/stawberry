package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zuzaaa-dev/stawberry/internal/app/apperror"
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/internal/domain/service/user"
	"github.com/zuzaaa-dev/stawberry/internal/handler/dto"
)

const basePath = ""

type UserService interface {
	CreateUser(ctx context.Context, user user.User) (string, string, error)
	Authenticate(ctx context.Context, email, password string) (string, error)
	GetUserByID(ctx context.Context, id string) (entity.User, error)
	UpdateUser(ctx context.Context, id string, updateUser user.UpdateUser) error
}

type userHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) userHandler {
	return userHandler{userService: userService}
}

func (h *userHandler) Registration(c *gin.Context) {
	var regUserDTO dto.RegistrationUserReq
	if err := c.ShouldBindJSON(&regUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apperror.BadRequest,
			"message": "Invalid user data",
			"details": err.Error(),
		})
		return
	}

	accessToken, refreshToken, err := h.userService.CreateUser(context.Background(), regUserDTO.ConvertToSvc())
	if err != nil {
		handleUserError(c, err)
		return
	}
	response := dto.RegistrationUserResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	jwtCookie := http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
		Path:  basePath + "/auth",
		// Domain:   "." + h.config.HTTPServer.Address,
		// MaxAge:   int(h.config.RefreshLife.Seconds()),
		Secure:   true,
		HttpOnly: true,
	}

	c.SetCookie(
		jwtCookie.Name,
		jwtCookie.Value,
		jwtCookie.MaxAge,
		jwtCookie.Path,
		jwtCookie.Domain,
		jwtCookie.Secure,
		jwtCookie.HttpOnly,
	)

	c.SetSameSite(http.SameSiteStrictMode)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var loginUserDto dto.LoginUserReq
	if err := c.ShouldBindJSON(&loginUserDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apperror.BadRequest,
			"message": "Invalid user data",
			"details": err.Error(),
		})
		return
	}

	response, err := h.userService.Authenticate(context.Background(), loginUserDto.Email, loginUserDto.Password)
	if err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
