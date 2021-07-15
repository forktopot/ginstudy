package common

//token

import (
	"ginstudy/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//创建一个jwt密钥
var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserId primitive.ObjectID
	jwt.StandardClaims
}

//设置token并发放
func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) //设置token过期时间
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(), //token发放的时间
			Issuer:    "forktopot.study", //谁发放的token
			Subject:   "user token",      //主题
		},
	}

	//使用jwtkey生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSting, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenSting, nil
}

//解析token并返回
func ParseToken(tokenSting string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenSting, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err
}
