package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"StickerDownloadBot/bot"
	"StickerDownloadBot/config"
)

var (
	version   = "dev"
	buildTime = "unknown"
	goVersion = "unknown"
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if ok {
		goVersion = info.GoVersion
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 8 {
				version = s.Value[:8]
				break
			}
		}
	}
	buildTime = time.Now().Format("2006-01-02 15:04:05")
}

func main() {
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		config.Setup()
	}

	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Println("读取config.yaml失败:", err)
		return
	}

	tgClient, qqClient, err := config.CreateClients(cfg.Proxy)
	if err != nil {
		fmt.Println("代理配置错误:", err)
		return
	}

	bot.New(cfg.Token, version, buildTime, goVersion, tgClient, qqClient).Start()
}
