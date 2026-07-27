package tools

import (
	"errors"
	"log"
	"strings"

	"github.com/golang-jwt/jwt"
)

// Sub2apiUser sub2api token 中需要的用户信息
type Sub2apiUser struct {
	UserID int64
	Email  string
	Role   string
}

// VerifySub2apiToken 纯离线校验 sub2api 签发的 JWT：HS256 签名 + exp/nbf 有效期。
// secret 与 sub2api 的 jwt.secret 保持一致（security_secrets 表 key=jwt_secret，
// 或部署机的 JWT_SECRET），配置项 Sub2apiJwtSecret / 环境变量 GOFLY_SUB2API_JWT_SECRET。
// 注意：离线校验不感知 sub2api 侧的吊销（改密码/退出登录），token 在自身过期前均有效。
func VerifySub2apiToken(tokenStr string) (*Sub2apiUser, error) {
	secret := GetServerConf().Sub2apiJwtSecret
	if secret == "" {
		return nil, errors.New("未配置 Sub2apiJwtSecret，无法校验 token")
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	user := &Sub2apiUser{}
	if v, ok := claims["user_id"].(float64); ok {
		user.UserID = int64(v)
	}
	if v, ok := claims["email"].(string); ok {
		user.Email = v
	}
	if v, ok := claims["role"].(string); ok {
		user.Role = v
	}
	if user.UserID == 0 {
		return nil, errors.New("token 缺少 user_id")
	}
	return user, nil
}

// LogSub2apiAuthStartupCheck 启动时检查离线验签配置
func LogSub2apiAuthStartupCheck() {
	if GetServerConf().Sub2apiJwtSecret == "" {
		message := "sub2api offline auth startup check failed: Sub2apiJwtSecret 未配置，带 token 的访客将无法通过验证"
		log.Println(message)
		Logger().Warn(message)
		return
	}
	message := "sub2api offline auth startup check ok: HS256 secret configured"
	log.Println(message)
	Logger().Info(message)
}

// IsEmail 简单的邮箱格式校验
func IsEmail(s string) bool {
	at := strings.LastIndex(s, "@")
	if at <= 0 || at == len(s)-1 {
		return false
	}
	domain := s[at+1:]
	dot := strings.LastIndex(domain, ".")
	return dot > 0 && dot < len(domain)-1 && !strings.ContainsAny(s, " \t\r\n")
}
