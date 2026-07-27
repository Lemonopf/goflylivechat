package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const FeishuWebhookBaseURL = "https://open.feishu.cn/open-apis/bot/v2/hook/"

func NormalizeFeishuWebhookToken(value string) string {
	token := strings.TrimSpace(value)
	if token == "" {
		return ""
	}
	if parsed, err := url.Parse(token); err == nil && parsed.Scheme != "" {
		token = parsed.Path
	} else {
		if index := strings.IndexAny(token, "?#"); index >= 0 {
			token = token[:index]
		}
	}
	token = strings.Trim(token, "/")
	if strings.Contains(token, "/") {
		parts := strings.Split(token, "/")
		token = parts[len(parts)-1]
	}
	return strings.TrimSpace(token)
}

func BuildFeishuWebhookURL(tokenOrURL string) string {
	token := NormalizeFeishuWebhookToken(tokenOrURL)
	if token == "" {
		return ""
	}
	return FeishuWebhookBaseURL + url.PathEscape(token)
}

func BuildFeishuTextPayload(text string) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": text,
		},
	})
}

func PostFeishuText(tokenOrURL string, text string) (string, error) {
	webhookURL := BuildFeishuWebhookURL(tokenOrURL)
	if webhookURL == "" {
		return "", errors.New("feishu webhook token is empty")
	}
	payload, err := BuildFeishuTextPayload(text)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
