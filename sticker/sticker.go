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
var idFromURLPattern = regexp.MustCompile(`[?&]id=(\d+)`)

// CleanName 清理文件名中的非法字符
func CleanName(name string) string {
	return unsafeNameChars.ReplaceAllString(name, "_")
}

// ExtractID 从用户输入中提取贴纸ID，支持直接ID或含 ?id= 的URL
func ExtractID(input string) string {
	if !strings.HasPrefix(input, "http") {
		return input
	}
	if match := idFromURLPattern.FindStringSubmatch(input); len(match) > 1 {
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败 (%d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	var data EmojiData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("解析失败")
	}
	return &data, nil
}

// Download 下载单个表情，文件名格式: 序号_名称.gif
func Download(client *http.Client, imgID, imgName string, height, index int, dir string) (string, error) {
	url := fmt.Sprintf("https://gxh.vip.qq.com/club/item/parcel/item/%s/%s/raw%d.gif", imgID[:2], imgID, height)
	filePath := filepath.Join(dir, fmt.Sprintf("%d_%s.gif", index, CleanName(imgName)))

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败 (%d)", resp.StatusCode)
	}

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
