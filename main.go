package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"

	"StickerDownloadBot/bot"
	"StickerDownloadBot/config"
)

var (
	version   = "dev"
	buildTime = "unknown"
	goVersion = runtime.Version()
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" && len(s.Value) >= 8 {
			version = s.Value[:8]
			break
		}
	}
}

func main() {
	if _, err := os.Stat("config.yaml"); errors.Is(err, os.ErrNotExist) {
		config.Setup()
	}

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal("读取config.yaml失败:", err)
	}

	tgClient, qqClient, err := config.CreateClients(cfg.Proxy)
	if err != nil {
		log.Fatal("代理配置错误:", err)
	}

	fmt.Printf("version: %s, buildTime: %s, goVersion: %s\n", version, buildTime, goVersion)
	bot.New(cfg.Token, version, buildTime, goVersion, tgClient, qqClient).Start()
}
