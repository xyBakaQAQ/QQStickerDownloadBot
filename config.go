package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
	"gopkg.in/yaml.v3"
)

type config struct {
	Token string      `yaml:"Token"`
	Proxy proxyConfig `yaml:"Proxy"`
}

type proxyConfig struct {
	Enable bool   `yaml:"Enable"`
	URL    string `yaml:"URL"`
}

func loadConfig(path string) (*config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("Token 未设置")
	}
	return &cfg, nil
}

func setupConfig() {
	fmt.Println("未检测到 config.yaml，开始初始化配置...")
	fmt.Println()

	data, err := os.ReadFile("config.example.yaml")
	if err != nil {
		fmt.Println("找不到 config.example.yaml，请手动创建 config.yaml")
		os.Exit(1)
	}
	os.WriteFile("config.yaml", data, 0644)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("请输入 Bot Token: ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)
	if token == "" {
		fmt.Println("Token 不能为空")
		os.Exit(1)
	}

	fmt.Print("是否启用代理? (y/n): ")
	proxyInput, _ := reader.ReadString('\n')
	proxyInput = strings.TrimSpace(strings.ToLower(proxyInput))

	enable := false
	proxyURL := ""
	if proxyInput == "y" || proxyInput == "yes" {
		fmt.Print("请输入代理地址 (如 socks5://127.0.0.1:7890): ")
		proxyURL, _ = reader.ReadString('\n')
		proxyURL = strings.TrimSpace(proxyURL)
		if proxyURL != "" {
			enable = true
		}
	}

	cfg := config{
		Token: token,
		Proxy: proxyConfig{Enable: enable, URL: proxyURL},
	}

	out, _ := yaml.Marshal(&cfg)
	os.WriteFile("config.yaml", out, 0644)

	fmt.Println()
	fmt.Println("配置已写入 config.yaml，启动中...")
	fmt.Println()
}

func createClients(cfg proxyConfig) (*http.Client, *http.Client, error) {
	qqClient := &http.Client{Timeout: 30 * time.Second}

	if !cfg.Enable || cfg.URL == "" {
		return &http.Client{Timeout: 30 * time.Second}, qqClient, nil
	}

	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("代理地址解析失败: %w", err)
	}

	transport := &http.Transport{}

	switch u.Scheme {
	case "http", "https":
		transport.Proxy = http.ProxyURL(u)
		fmt.Println("使用HTTP代理:", cfg.URL)
	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", u.Host, nil, proxy.Direct)
		if err != nil {
			return nil, nil, fmt.Errorf("SOCKS5代理配置失败: %w", err)
		}
		transport.DialContext = dialer.(proxy.ContextDialer).DialContext
		fmt.Println("使用SOCKS5代理:", cfg.URL)
	default:
		return nil, nil, fmt.Errorf("不支持的代理协议: %s（支持 http / socks5）", u.Scheme)
	}

	return &http.Client{Transport: transport, Timeout: 30 * time.Second}, qqClient, nil
}
