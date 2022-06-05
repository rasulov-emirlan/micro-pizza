package jwtlib

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/internal/domain"
)

type jwtmanager struct {
	key        []byte
	accessExp  time.Duration
	refreshExp time.Duration
}

func NewJwtManager(key []byte) *jwtmanager {
	return &jwtmanager{key: key}
}

func (j *jwtmanager) SetKey(key []byte) {
	j.key = key
}

func (j *jwtmanager) SetExp(accessExp, refreshExp time.Duration) {
	j.accessExp = accessExp
	j.refreshExp = refreshExp
}

type Claims struct {
	jwt.StandardClaims
	UserID domain.ID
	Roles  []domain.Role
}

type RefreshClaims struct {
	jwt.StandardClaims
	UserID domain.ID
}

func (j *jwtmanager) Generate(userID domain.ID, roles []domain.Role) (domain.SignInOutput, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.accessExp).Unix(),
		},
		UserID: userID,
		Roles:  roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(j.key)
	if err != nil {
		return domain.SignInOutput{}, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.refreshExp).Unix(),
		},
	})
	refreshKey, err := refreshToken.SignedString(j.key)
	if err != nil {
		return domain.SignInOutput{}, err
	}

	return domain.SignInOutput{
		AccessKey:  accessToken,
		RefreshKey: refreshKey,
	}, nil
}

func (j *jwtmanager) DecodeAccess(accessKey string) (domain.AccessClaims, error) {
	token, err := jwt.ParseWithClaims(accessKey, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})
	if err != nil {
		return domain.AccessClaims{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return domain.AccessClaims{}, err
	}

	return domain.AccessClaims{
		UserID: claims.UserID,
		Roles:  claims.Roles,
	}, nil
}

func (j *jwtmanager) DecodeRefresh(refreshKey string) (domain.RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(refreshKey, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})
	if err != nil {
		return domain.RefreshClaims{}, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok {
		return domain.RefreshClaims{}, err
	}

	return domain.RefreshClaims{
		UserID: claims.UserID,
	}, nil
}
