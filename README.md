# StickerDownloadBot

> **请不要将 Bot Token 上传到公开仓库**

Telegram Bot，下载 QQ 表情贴纸包。发送链接或 ID，自动下载全部 GIF 打包为 zip 发回。


## 下载

从 [Releases](https://github.com/xyBakaQAQ/QQStickerDownloadBot/releases/tag/latest) 页面下载对应平台的可执行文件。

## 构建

```bash
git clone https://github.com/xyBakaQAQ/StickerDownloadBot.git
cd StickerDownloadBot

# Windows
build.bat

# Linux / macOS
./build.sh
```

构建脚本会自动注入版本号和构建时间，`/about` 命令可以查看。

也可以手动构建：

```bash
go build -trimpath -ldflags="-s -w -X 'main.buildTime=$(date '+%Y-%m-%d %H:%M:%S')'" .
```

## 配置

在 [@BotFather](https://t.me/BotFather) 创建 Bot，获取 Token

```yaml
Token: "你的Token"

Proxy:
  Enable: true
  URL: "socks5://127.0.0.1:7890"
```

> 代理仅用于 Telegram API，贴纸下载直连 QQ 服务器。

## 依赖

- Go 1.25 及以上版本

## 许可证

[**MIT License**](LICENSE)