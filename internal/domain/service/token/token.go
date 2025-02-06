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

const expires = time.Hour * 24 * 30

type tokenService struct {
	jwtSecret string
}

func NewTokenService(secret string) *tokenService {
	return &tokenService{jwtSecret: secret}
}

func (ts *tokenService) GenerateTokens(ctx context.Context, fingerprint string, userID uint) (string, entity.RefreshToken, error) {
	accessToken, err := generateJWT(userID, ts.jwtSecret, time.Hour*1)
	if err != nil {
		return "", entity.RefreshToken{}, err
	}

	entityRefreshToken, err := generateRefresh(fingerprint, userID)
	if err != nil {
		return "", entity.RefreshToken{}, err
	}

	return accessToken, entityRefreshToken, nil
}

func (ts *tokenService) Parse(token string) (entity.AccessToken, error) {
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

func generateJWT(userID uint, secret string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func generateRefresh(fingerprint string, userID uint) (entity.RefreshToken, error) {
	now := time.Now()

	refreshUUID, err := uuid.NewRandom()
	if err != nil {
		return entity.RefreshToken{}, err
	}

	return entity.RefreshToken{
		UUID:        refreshUUID,
		CreatedAt:   now,
		ExpiresAt:   now.Add(expires),
		RevokedAt:   nil,
		Fingerprint: fingerprint,
		UserID:      userID,
	}, nil
}
