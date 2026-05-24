package sticker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// EmojiData QQ表情包元数据
type EmojiData struct {
	SupportSize []struct{ Width, Height int } `json:"supportSize"`
	Name        string                        `json:"name"`
	Mark        string                        `json:"mark"`
	Imgs        []struct{ ID, Name string }   `json:"imgs"`
}

var unsafeNameChars = regexp.MustCompile(`[<>:"/\\|?*]`)

// CleanName 清理文件名中的非法字符
func CleanName(name string) string {
	return unsafeNameChars.ReplaceAllString(name, "_")
}

// ExtractID 从用户输入中提取贴纸ID，支持直接ID或含 ?id= 的URL
func ExtractID(input string) string {
	if !strings.HasPrefix(input, "http") {
		return input
	}
	if match := regexp.MustCompile(`[?&]id=(\d+)`).FindStringSubmatch(input); len(match) > 1 {
		return match[1]
	}
	return ""
}

// FetchMeta 从QQ服务器获取表情包元数据
func FetchMeta(client *http.Client, id string) (*EmojiData, error) {
	url := fmt.Sprintf("https://gxh.vip.qq.com/club/item/parcel/%s/%s.json", id[len(id)-1:], id)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data EmojiData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析失败")
	}
	return &data, nil
}

// DownloadOne 下载单个表情，文件名格式: 序号_名称_ID.gif
func DownloadOne(client *http.Client, imgID, imgName string, height, index int, dir string) (string, error) {
	url := fmt.Sprintf("https://gxh.vip.qq.com/club/item/parcel/item/%s/%s/raw%d.gif", imgID[:2], imgID, height)
	filePath := filepath.Join(dir, fmt.Sprintf("%d_%s.gif", index, CleanName(imgName)))

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", err
	}
	return filePath, nil
}
