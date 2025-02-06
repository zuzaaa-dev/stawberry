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
	CreateUser(ctx context.Context, user user.User, fingerprint string) (string, string, error)
	Authenticate(ctx context.Context, email, password, fingerprint string) (string, string, error)
	Refresh(ctx context.Context, refreshToken, fingerprint string) (string, string, error)
	Logout(ctx context.Context, refreshToken, fingerprint string) error
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

	accessToken, refreshToken, err := h.userService.CreateUser(
		context.Background(),
		regUserDTO.ConvertToSvc(),
		regUserDTO.Fingerprint,
	)
	if err != nil {
		handleUserError(c, err)
		return
	}
	response := dto.RegistrationUserResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	setRefreshCookie(c, refreshToken, "", 1)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var loginUserDTO dto.LoginUserReq
	if err := c.ShouldBindJSON(&loginUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apperror.BadRequest,
			"message": "Invalid user data",
			"details": err.Error(),
		})
		return
	}

	accessToken, refreshToken, err := h.userService.Authenticate(
		context.Background(),
		loginUserDTO.Email,
		loginUserDTO.Password,
		loginUserDTO.Fingerprint,
	)
	if err != nil {
		handleUserError(c, err)
		return
	}

	response := dto.LoginUserResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	setRefreshCookie(c, refreshToken, "", 1)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Refresh(c *gin.Context) {
	var refreshDTO dto.RefreshReq
	if err := c.ShouldBindJSON(&refreshDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apperror.BadRequest,
			"message": "Invalid refresh data",
			"details": err.Error(),
		})
		return
	}

	if refreshDTO.RefreshToken == "" {
		refresh, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    apperror.BadRequest,
				"message": "Invalid refresh data",
				"details": err.Error(),
			})
			return
		}
		refreshDTO.RefreshToken = refresh
	}

	accessToken, refreshToken, err := h.userService.Refresh(
		context.Background(),
		refreshDTO.RefreshToken,
		refreshDTO.Fingerprint,
	)
	if err != nil {
		handleUserError(c, err)
		return
	}

	response := dto.RefreshResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	setRefreshCookie(c, refreshToken, "", 0)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Logout(c *gin.Context) {
	var logoutDTO dto.LogoutReq
	if err := c.ShouldBindJSON(&logoutDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    apperror.BadRequest,
			"message": "Invalid refresh data",
			"details": err.Error(),
		})
		return
	}

	if logoutDTO.RefreshToken == "" {
		refresh, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    apperror.BadRequest,
				"message": "Invalid refresh data",
				"details": err.Error(),
			})
			return
		}
		logoutDTO.RefreshToken = refresh
	}

	if err := h.userService.Logout(
		context.Background(),
		logoutDTO.RefreshToken,
		logoutDTO.Fingerprint,
	); err != nil {
		handleUserError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func setRefreshCookie(c *gin.Context, refreshToken, domain string, maxAge int) {
	jwtCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     basePath + "/auth",
		Domain:   domain,
		MaxAge:   maxAge,
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
}
