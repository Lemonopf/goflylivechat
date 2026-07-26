package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	Sub2apiAuthReasonUnknown              = "SUB2API_AUTH_UNKNOWN"
	Sub2apiAuthReasonConfigMissing        = "SUB2API_API_URL_MISSING"
	Sub2apiAuthReasonInvalidBaseURL       = "SUB2API_API_URL_INVALID"
	Sub2apiAuthReasonRequestBuildFailed   = "SUB2API_REQUEST_BUILD_FAILED"
	Sub2apiAuthReasonNetworkError         = "SUB2API_NETWORK_ERROR"
	Sub2apiAuthReasonRejected             = "SUB2API_AUTH_REJECTED"
	Sub2apiAuthReasonResponseDecodeFailed = "SUB2API_RESPONSE_DECODE_FAILED"
	Sub2apiAuthReasonResponseInvalid      = "SUB2API_RESPONSE_INVALID"
)

// Sub2apiUser sub2api 在线验证返回的用户信息
type Sub2apiUser struct {
	UserID int64
	Email  string
	Role   string
}

// Sub2apiAuthError 保留内部诊断原因，避免把 token 或上游细节暴露给访客端。
type Sub2apiAuthError struct {
	Reason     string
	StatusCode int
	Message    string
	Err        error
}

func (e *Sub2apiAuthError) Error() string {
	parts := []string{e.Reason}
	if e.StatusCode != 0 {
		parts = append(parts, fmt.Sprintf("http=%d", e.StatusCode))
	}
	if e.Message != "" {
		parts = append(parts, "message="+e.Message)
	}
	if e.Err != nil {
		parts = append(parts, "err="+e.Err.Error())
	}
	return strings.Join(parts, " ")
}

func (e *Sub2apiAuthError) Unwrap() error {
	return e.Err
}

func Sub2apiAuthReason(err error) string {
	var authErr *Sub2apiAuthError
	if errors.As(err, &authErr) && authErr.Reason != "" {
		return authErr.Reason
	}
	return Sub2apiAuthReasonUnknown
}

// CheckSub2apiAuthConfig 校验在线验签所需配置是否可用于构造 /api/v1/auth/me 地址。
func CheckSub2apiAuthConfig() error {
	_, err := sub2apiAuthMeURL(GetServerConf().Sub2apiApiURL)
	return err
}

func LogSub2apiAuthStartupCheck() {
	baseURL := strings.TrimSpace(GetServerConf().Sub2apiApiURL)
	endpoint, err := sub2apiAuthMeURL(baseURL)
	if err != nil {
		message := fmt.Sprintf("sub2api online auth startup check failed: reason=%s base_url=%q err=%v", Sub2apiAuthReason(err), baseURL, err)
		log.Println(message)
		Logger().Warn(message)
		return
	}
	message := fmt.Sprintf("sub2api online auth startup check ok: base_url=%q endpoint=%q", strings.TrimRight(baseURL, "/"), endpoint)
	log.Println(message)
	Logger().Info(message)
}

// VerifySub2apiToken 在线校验 sub2api 签发的 token：
// 调用 sub2api 的 GET /api/v1/auth/me（Authorization: Bearer），由 sub2api 完成签名、有效期、吊销等全部校验。
// 客服系统无需保存 sub2api 的 jwt.secret。
// clientIP / userAgent 为访客的真实 IP 和 UA，需原样转发以通过 sub2api 的会话绑定（IP+UA 指纹）校验，
// 否则 sub2api 会判定指纹不匹配并吊销该用户会话。
// 配置项 Sub2apiApiURL / 环境变量 GOFLY_SUB2API_API_URL，例如 https://lemonzz.xyz
func VerifySub2apiToken(token, clientIP, userAgent string) (*Sub2apiUser, error) {
	authURL, err := sub2apiAuthMeURL(GetServerConf().Sub2apiApiURL)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, authURL, nil)
	if err != nil {
		return nil, &Sub2apiAuthError{Reason: Sub2apiAuthReasonRequestBuildFailed, Err: err}
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
		return nil, &Sub2apiAuthError{Reason: Sub2apiAuthReasonNetworkError, Err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, &Sub2apiAuthError{
			Reason:     Sub2apiAuthReasonRejected,
			StatusCode: resp.StatusCode,
			Message:    readSub2apiErrorMessage(resp.Body),
		}
	}
	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Msg     string `json:"msg"`
		Data    struct {
			ID    int64  `json:"id"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, &Sub2apiAuthError{Reason: Sub2apiAuthReasonResponseDecodeFailed, Err: err}
	}
	if result.Code != 0 {
		return nil, &Sub2apiAuthError{
			Reason:  Sub2apiAuthReasonRejected,
			Message: firstNonEmpty(result.Message, result.Msg),
		}
	}
	if result.Data.ID == 0 {
		return nil, &Sub2apiAuthError{Reason: Sub2apiAuthReasonResponseInvalid, Message: "missing user id"}
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

func sub2apiAuthMeURL(baseURL string) (string, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return "", &Sub2apiAuthError{Reason: Sub2apiAuthReasonConfigMissing, Message: "missing Sub2apiApiURL"}
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", &Sub2apiAuthError{Reason: Sub2apiAuthReasonInvalidBaseURL, Message: baseURL, Err: err}
	}
	return baseURL + "/api/v1/auth/me", nil
}

func readSub2apiErrorMessage(body io.Reader) string {
	var result struct {
		Code    int    `json:"code"`
		Reason  string `json:"reason"`
		Message string `json:"message"`
		Msg     string `json:"msg"`
		Error   string `json:"error"`
	}
	if err := json.NewDecoder(io.LimitReader(body, 4096)).Decode(&result); err != nil {
		return ""
	}
	message := firstNonEmpty(result.Reason, result.Message, result.Msg, result.Error)
	if result.Code != 0 && message != "" {
		return fmt.Sprintf("code=%d %s", result.Code, message)
	}
	if result.Code != 0 {
		return fmt.Sprintf("code=%d", result.Code)
	}
	return message
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
