package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const wechatAPI = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

type Authenticator struct {
	appID     string
	secret    string
	jwtSecret []byte
}

func NewAuthenticator() *Authenticator {
	return &Authenticator{
		appID:     os.Getenv("WECHAT_APPID"),
		secret:    os.Getenv("WECHAT_APPSECRET"),
		jwtSecret: []byte(os.Getenv("JWT_SECRET")),
	}
}

// Login exchanges a WeChat code for a user session and returns a JWT.
func (a *Authenticator) Login(code string) (string, error) {
	url := fmt.Sprintf(wechatAPI, a.appID, a.secret, code)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var session struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return "", err
	}

	if session.ErrCode != 0 {
		return "", fmt.Errorf("wechat login failed: %s", session.ErrMsg)
	}

	// Create a new token object, specifying signing method and the claims
	claims := jwt.MapClaims{
		"sub": session.OpenID, // Subject (user identifier)
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Token expires in 30 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT and returns the userID (openid) from its claims.
func (a *Authenticator) ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["sub"].(string); ok {
			return userID, nil
		}
	}

	return "", fmt.Errorf("invalid token")
}
