package jwt

import (
	"errors"
	"gin-mall/consts"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type EmailClaims struct {
	UserId        uint   `json:"user_id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OperationType uint   `json:"operation_type"`
	jwt.StandardClaims
}

var jwtSecret = []byte("takatsuki")

func GenerateToken(id uint, username string) (accessToken string, refreshToken string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(consts.AccessTokenExpireDuration)
	rtExpireTime := nowTime.Add(consts.RefreshTokenExpireDuration)
	claims := Claims{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "mall",
		},
	}
	reclaims := Claims{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: rtExpireTime.Unix(),
			Issuer:    "mall",
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, reclaims).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func ParseRefreshToken(aToken, rToken string) (newAtoken, newRtoken string, err error) {
	accessClaim, err := ParseToken(aToken)
	if err != nil {
		return
	}
	refreshClaim, err := ParseToken(rToken)
	if err != nil {
		return
	}
	//access过期，重新生成
	if accessClaim.ExpiresAt > time.Now().Unix() {
		return GenerateToken(accessClaim.ID, accessClaim.Username)
	}
	//refresh过期，重新生成
	if refreshClaim.ExpiresAt > time.Now().Unix() {
		return GenerateToken(refreshClaim.ID, refreshClaim.Username)
	}
	return "", "", errors.New("身份过期，请重新登录")
}

// 签发邮箱验证token
func GenerateEmailToken(userID, operation uint, email, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(15 * time.Minute)
	claims := EmailClaims{
		UserId:        userID,
		Email:         email,
		Password:      password,
		OperationType: operation,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "mall",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

// 解析邮件token
func ParseEmailToken(token string) (*EmailClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &EmailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*EmailClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
