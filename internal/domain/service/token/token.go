package token

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zuzaaa-dev/stawberry/internal/app/apperror"

	"github.com/golang-jwt/jwt"
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
)

var signingMethod = jwt.SigningMethodHS256

type Repository interface {
	InsertToken(ctx context.Context, token entity.RefreshToken) error
	GetActivesTokenByUserID(ctx context.Context, userID uint) ([]entity.RefreshToken, error)
	RevokeActivesByUserID(ctx context.Context, userID uint) error
	GetByUUID(ctx context.Context, uuid string) (entity.RefreshToken, error)
	Update(ctx context.Context, refresh entity.RefreshToken) (entity.RefreshToken, error)
}

type tokenService struct {
	tokenRepository Repository
	jwtSecret       string
	refreshLife     time.Duration
	accessLife      time.Duration
}

func NewTokenService(tokenRepo Repository, secret string) *tokenService {
	return &tokenService{
		tokenRepository: tokenRepo,
		jwtSecret:       secret,
	}
}

func (ts *tokenService) GenerateTokens(
	ctx context.Context,
	fingerprint string,
	userID uint,
) (string, entity.RefreshToken, error) {
	accessToken, err := generateJWT(userID, ts.jwtSecret, ts.accessLife)
	if err != nil {
		return "", entity.RefreshToken{}, err
	}

	entityRefreshToken, err := generateRefresh(fingerprint, userID, ts.refreshLife)
	if err != nil {
		return "", entity.RefreshToken{}, err
	}

	return accessToken, entityRefreshToken, nil
}

func (ts *tokenService) ValidateToken(
	ctx context.Context,
	token string,
) (entity.AccessToken, error) {
	accessToken, err := ts.parse(token)
	if err != nil {
		return entity.AccessToken{}, err
	}

	if time.Now().After(accessToken.ExpiresAt) {
		return entity.AccessToken{}, apperror.ErrInvalidToken
	}

	return accessToken, nil
}

func (ts *tokenService) parse(token string) (entity.AccessToken, error) {
	claim := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		if token.Header["alg"] != signingMethod.Alg() {
			appError := apperror.ErrInvalidToken
			appError.Err = fmt.Errorf("invalid signing method")
			return nil, appError
		}
		return []byte(ts.jwtSecret), nil
	})
	if err != nil {
		return entity.AccessToken{}, apperror.ErrInvalidToken
	}
	userID, ok := claim["sub"].(float64)
	if !ok {
		return entity.AccessToken{}, apperror.ErrInvalidToken
	}

	unixExpiresAt, ok := claim["exp"].(float64)
	if !ok {
		return entity.AccessToken{}, apperror.ErrInvalidToken
	}
	expiresAt := time.Unix(int64(unixExpiresAt), 0)

	unixIssuedAt, ok := claim["iat"].(float64)
	if !ok {
		return entity.AccessToken{}, apperror.ErrInvalidToken
	}

	issuedAt := time.Unix(int64(unixIssuedAt), 0)

	return entity.AccessToken{
		UserID:    uint(userID),
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}, nil
}

func (ts *tokenService) InsertToken(
	ctx context.Context,
	token entity.RefreshToken,
) error {
	return ts.tokenRepository.InsertToken(ctx, token)
}

func (ts *tokenService) GetActivesTokenByUserID(
	ctx context.Context,
	userID uint,
) ([]entity.RefreshToken, error) {
	return ts.tokenRepository.GetActivesTokenByUserID(ctx, userID)
}

func (ts *tokenService) RevokeActivesByUserID(
	ctx context.Context,
	userID uint,
) error {
	return ts.tokenRepository.RevokeActivesByUserID(ctx, userID)
}

func (ts *tokenService) GetByUUID(
	ctx context.Context,
	uuid string,
) (entity.RefreshToken, error) {
	return ts.tokenRepository.GetByUUID(ctx, uuid)
}

func (ts *tokenService) Update(
	ctx context.Context,
	refresh entity.RefreshToken,
) (entity.RefreshToken, error) {
	return ts.tokenRepository.Update(ctx, refresh)
}

func generateJWT(userID uint, secret string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func generateRefresh(fingerprint string, userID uint, refreshLife time.Duration) (entity.RefreshToken, error) {
	now := time.Now()

	refreshUUID, err := uuid.NewRandom()
	if err != nil {
		return entity.RefreshToken{}, err
	}

	return entity.RefreshToken{
		UUID:        refreshUUID,
		CreatedAt:   now,
		ExpiresAt:   now.Add(refreshLife),
		RevokedAt:   nil,
		Fingerprint: fingerprint,
		UserID:      userID,
	}, nil
}
