# StickerDownloadBot

> **请不要将 Bot Token 上传到公开仓库**

Telegram Bot，下载 QQ 表情贴纸包。发送链接或 ID，自动下载全部 GIF 打包为 zip 发回。


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

- Go 1.24 及以上版本

## 许可证

**MIT License**