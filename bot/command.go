package bot

import (
	"context"
	"fmt"
	"strings"

	lib "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) dispatchCommand(ctx context.Context, api *lib.Bot, update *models.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	input := strings.TrimSpace(update.Message.Text)
	chatID := update.Message.Chat.ID

	switch {
	case input == "/start":
		b.send(ctx, chatID, "发送贴纸链接或ID即可下载\n使用 /help 查看帮助")
	case input == "/help":
		b.send(ctx, chatID, "发送贴纸链接或ID即可下载全部表情\n\n怎么获取链接:\n手机QQ：\n「表情详情」页面 → 右上角 → 复制链接\n\nQQNT:\nEmoji\\marketface\\json 目录下\n文件名以 ID 开头的 jtmp 文件\n把文件名里的数字 ID 发给 Bot 即可\n\n命令:\n/start - 开始\n/help - 帮助\n/about - 关于")
	case input == "/about":
		b.send(ctx, chatID, fmt.Sprintf("版本: %s\nGo: %s\n构建时间: %s\n用户ID: %d", b.version, b.goVersion, b.buildTime, chatID))
	case strings.HasPrefix(input, "/"):
		return true
	default:
		return false
	}
	return true
}
