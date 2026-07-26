package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Sub2apiUser sub2api 在线验证返回的用户信息
type Sub2apiUser struct {
	UserID int64
	Email  string
	Role   string
}

// VerifySub2apiToken 在线校验 sub2api 签发的 token：
// 调用 sub2api 的 GET /api/v1/auth/me（Authorization: Bearer），由 sub2api 完成签名、有效期、吊销等全部校验。
// 客服系统无需保存 sub2api 的 jwt.secret。
// clientIP / userAgent 为访客的真实 IP 和 UA，需原样转发以通过 sub2api 的会话绑定（IP+UA 指纹）校验，
// 否则 sub2api 会判定指纹不匹配并吊销该用户会话。
// 配置项 Sub2apiApiURL / 环境变量 GOFLY_SUB2API_API_URL，例如 https://lemonzz.xyz
func VerifySub2apiToken(token, clientIP, userAgent string) (*Sub2apiUser, error) {
	baseURL := strings.TrimRight(GetServerConf().Sub2apiApiURL, "/")
	if baseURL == "" {
		return nil, errors.New("未配置 Sub2apiApiURL，无法校验 token")
	}
	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/v1/auth/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
	if clientIP != "" {
		req.Header.Set("X-Real-IP", clientIP)
		req.Header.Set("X-Forwarded-For", clientIP)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("验证服务暂不可用: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token 无效或已过期（HTTP %d）", resp.StatusCode)
	}
	var result struct {
		Code int `json:"code"`
		Data struct {
			ID    int64  `json:"id"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析验证响应失败: %v", err)
	}
	if result.Data.ID == 0 {
		return nil, errors.New("token 无效")
	}
	return &Sub2apiUser{
		UserID: result.Data.ID,
		Email:  result.Data.Email,
		Role:   result.Data.Role,
	}, nil
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
