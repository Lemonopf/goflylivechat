package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goflylivechat/models"
	"goflylivechat/tools"
	"log"
	"strconv"
	"strings"
)

const (
	configNoticeFeishuWebhook = "NoticeFeishuWebhook"
	configFeishuWebhookToken  = "FeishuWebhookToken"
)

func SendFeishuWebhookNotice(kefuName, visitorName, content, visitorId, mainURL string) {
	enabled, err := strconv.ParseBool(models.FindConfigByUserId(kefuName, configNoticeFeishuWebhook).ConfValue)
	if err != nil || !enabled {
		return
	}
	token := models.FindConfigByUserId(kefuName, configFeishuWebhookToken).ConfValue
	if tools.NormalizeFeishuWebhookToken(token) == "" {
		return
	}
	text := buildFeishuNoticeText(visitorName, content, visitorId, mainURL)
	res, err := tools.PostFeishuText(token, text)
	if err != nil {
		log.Println("send feishu webhook error:", err)
		return
	}
	tools.Logger().Infoln("send_feishu_webhook", kefuName, res)
}

func buildFeishuNoticeText(visitorName, content, visitorId, mainURL string) string {
	lines := []string{
		"新客户消息",
		fmt.Sprintf("访客：%s", visitorName),
		fmt.Sprintf("内容：%s", content),
	}
	if visitorId != "" {
		lines = append(lines, fmt.Sprintf("访客ID：%s", visitorId))
	}
	if mainURL != "" {
		lines = append(lines, fmt.Sprintf("客服后台：%s", mainURL))
	}
	return strings.Join(lines, "\n")
}

func mainURLFromRequest(c *gin.Context) string {
	if c == nil || c.Request == nil || c.Request.Host == "" {
		return ""
	}
	host := firstForwardedHeaderValue(c.GetHeader("X-Forwarded-Host"))
	if host == "" {
		host = c.Request.Host
	}
	scheme := c.GetHeader("X-Forwarded-Proto")
	scheme = firstForwardedHeaderValue(scheme)
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	return scheme + "://" + host + "/main"
}

func firstForwardedHeaderValue(value string) string {
	items := strings.Split(value, ",")
	if len(items) == 0 {
		return ""
	}
	return strings.TrimSpace(items[0])
}
