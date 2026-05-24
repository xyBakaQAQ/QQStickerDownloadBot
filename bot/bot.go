package bot

import (
	"archive/zip"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"StickerDownloadBot/sticker"

	lib "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Bot Telegram机器人
type Bot struct {
	api       *lib.Bot
	qqClient  *http.Client
	version   string
	buildTime string
	goVersion string
}

// New 创建Bot
// tgClient 走代理，用于Telegram API
// qqClient 直连，用于QQ贴纸下载
func New(token, version, buildTime, goVersion string, tgClient, qqClient *http.Client) *Bot {
	b := &Bot{qqClient: qqClient, version: version, buildTime: buildTime, goVersion: goVersion}

	api, err := lib.New(token,
		lib.WithHTTPClient(30*time.Second, tgClient),
		lib.WithDefaultHandler(b.handle),
	)
	if err != nil {
		fmt.Println("Bot初始化失败:", err)
		return nil
	}

	b.api = api
	return b
}

// Start 启动Bot
func (b *Bot) Start() {
	user, err := b.api.GetMe(context.TODO())
	if err != nil {
		fmt.Println("获取Bot信息失败:", err)
		return
	}
	fmt.Printf("Bot启动成功: @%s\n", user.Username)
	b.api.Start(context.TODO())
}

func (b *Bot) handle(ctx context.Context, api *lib.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	chat := update.Message.Chat
	from := update.Message.From
	who := fmt.Sprintf("%d", chat.ID)
	if from != nil {
		who = from.FirstName
		if from.Username != "" {
			who += "(@" + from.Username + ")"
		}
	}
	fmt.Printf("[%s] %s\n", who, update.Message.Text)

	if b.dispatchCommand(ctx, api, update) {
		return
	}

	input := strings.TrimSpace(update.Message.Text)
	chatID := chat.ID

	id := sticker.ExtractID(input)
	if id == "" {
		b.send(ctx, chatID, "无效的链接或ID")
		return
	}

	b.processRemote(ctx, chatID, id)
}

func (b *Bot) send(ctx context.Context, chatID int64, text string) *models.Message {
	msg, err := b.api.SendMessage(ctx, &lib.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		fmt.Println("send error:", err)
	}
	return msg
}

func (b *Bot) edit(ctx context.Context, chatID int64, msgID int, text string) {
	b.api.EditMessageText(ctx, &lib.EditMessageTextParams{
		ChatID:    chatID,
		MessageID: msgID,
		Text:      text,
	})
}

func (b *Bot) processRemote(ctx context.Context, chatID int64, id string) {
	data, err := sticker.FetchMeta(b.qqClient, id)
	if err != nil {
		b.send(ctx, chatID, "请求失败: "+err.Error())
		return
	}

	dir, err := os.MkdirTemp("", "sticker_"+id)
	if err != nil {
		b.send(ctx, chatID, "创建临时目录失败")
		return
	}
	defer os.RemoveAll(dir)

	info := fmt.Sprintf("名称: %s\nID: %s\n描述: %s", data.Name, id, data.Mark)
	total := len(data.Imgs)

	status := b.send(ctx, chatID, fmt.Sprintf("%s\n表情数: %d\n[0/%d] 下载中...", info, total, total))
	if status == nil {
		return
	}

	h := data.SupportSize[0].Height
	success := 0
	failed := 0

	for i, img := range data.Imgs {
		_, err := sticker.DownloadOne(b.qqClient, img.ID, img.Name, h, i+1, dir)
		if err != nil {
			failed++
		} else {
			success++
		}

		b.edit(ctx, chatID, status.ID,
			fmt.Sprintf("%s\n表情数: %d\n[%d/%d] 下载中...", info, total, i+1, total))

		time.Sleep(100 * time.Millisecond)
	}

	if success == 0 {
		b.edit(ctx, chatID, status.ID, info+"\n所有表情下载失败")
		return
	}

	b.edit(ctx, chatID, status.ID, fmt.Sprintf("%s\n表情数: %d\n下载完成，正在压缩...", info, success))

	zipName := fmt.Sprintf("[ID%s] %s.zip", id, data.Name)
	zipPath := filepath.Join(dir, zipName)
	if err := zipDir(dir, zipPath); err != nil {
		b.edit(ctx, chatID, status.ID, fmt.Sprintf("%s\n压缩失败: %s", info, err.Error()))
		return
	}

	file, err := os.Open(zipPath)
	if err != nil {
		b.edit(ctx, chatID, status.ID, "读取压缩包失败")
		return
	}
	defer file.Close()

	b.api.SendDocument(ctx, &lib.SendDocumentParams{
		ChatID: chatID,
		Document: &models.InputFileUpload{
			Filename: zipName,
			Data:     file,
		},
	})

	if failed > 0 {
		b.edit(ctx, chatID, status.ID, fmt.Sprintf("%s\n共 %d 个，下载完成（%d 个失败）", info, success, failed))
	} else {
		b.edit(ctx, chatID, status.ID, fmt.Sprintf("%s\n共 %d 个，下载完成", info, success))
	}
}

func zipDir(srcDir, zipPath string) error {
	f, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	// 按文件名前缀数字排序: 1_xxx.gif < 2_xxx.gif < 10_xxx.gif
	sort.Slice(entries, func(i, j int) bool {
		ni, _ := strconv.Atoi(strings.SplitN(entries[i].Name(), "_", 2)[0])
		nj, _ := strconv.Atoi(strings.SplitN(entries[j].Name(), "_", 2)[0])
		return ni < nj
	})

	skipName := filepath.Base(zipPath)
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == skipName {
			continue
		}
		data, err := os.ReadFile(filepath.Join(srcDir, entry.Name()))
		if err != nil {
			continue
		}
		writer, err := w.Create(entry.Name())
		if err != nil {
			continue
		}
		writer.Write(data)
	}
	return nil
}
