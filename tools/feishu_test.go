package tools

import (
	"encoding/json"
	"testing"
)

func TestNormalizeFeishuWebhookToken(t *testing.T) {
	token := "test-token-123"
	tests := map[string]string{
		token:               token,
		"  " + token + "  ": token,
		"https://open.feishu.cn/open-apis/bot/v2/hook/" + token:                token,
		"https://open.feishu.cn/open-apis/bot/v2/hook/" + token + "?foo=bar":   token,
		"https://open.larksuite.com/open-apis/bot/v2/hook/" + token + "#debug": token,
	}

	for input, want := range tests {
		if got := NormalizeFeishuWebhookToken(input); got != want {
			t.Fatalf("NormalizeFeishuWebhookToken(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestBuildFeishuWebhookURL(t *testing.T) {
	token := "test-token-123"
	want := "https://open.feishu.cn/open-apis/bot/v2/hook/" + token
	if got := BuildFeishuWebhookURL(token); got != want {
		t.Fatalf("BuildFeishuWebhookURL(%q) = %q, want %q", token, got, want)
	}
}

func TestBuildFeishuTextPayload(t *testing.T) {
	payload, err := BuildFeishuTextPayload("新客户消息")
	if err != nil {
		t.Fatalf("BuildFeishuTextPayload returned error: %v", err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(payload, &got); err != nil {
		t.Fatalf("payload is not valid JSON: %v", err)
	}
	if got["msg_type"] != "text" {
		t.Fatalf("msg_type = %v, want text", got["msg_type"])
	}
	content := got["content"].(map[string]interface{})
	if content["text"] != "新客户消息" {
		t.Fatalf("content.text = %v, want 新客户消息", content["text"])
	}
}
