package config

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Token string      `yaml:"Token"`
	Proxy ProxyConfig `yaml:"Proxy"`
}

type ProxyConfig struct {
	Enable bool   `yaml:"Enable"`
	URL    string `yaml:"URL"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("Token 未设置")
	}
	return &cfg, nil
}

func Setup() {
	fmt.Println("未检测到 config.yaml，开始初始化配置...\n")

	var token, proxyInput, proxyURL string
	fmt.Print("请输入 Bot Token: ")
	fmt.Scanln(&token)
	if token == "" {
		fmt.Println("Token 不能为空")
		os.Exit(1)
	}

	fmt.Print("是否启用代理? (y/n): ")
	fmt.Scanln(&proxyInput)

	enable := false
	if strings.HasPrefix(strings.ToLower(proxyInput), "y") {
		fmt.Print("请输入代理地址 (如 socks5://127.0.0.1:7890): ")
		fmt.Scanln(&proxyURL)
		enable = proxyURL != ""
	}

	cfg := Config{Token: token, Proxy: ProxyConfig{Enable: enable, URL: proxyURL}}
	out, _ := yaml.Marshal(&cfg)
	os.WriteFile("config.yaml", out, 0644)

	fmt.Println("\n配置已写入 config.yaml，启动中...\n")
}

func newClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}

func CreateClients(cfg ProxyConfig) (*http.Client, *http.Client, error) {
	if !cfg.Enable || cfg.URL == "" {
		return newClient(30 * time.Second), newClient(30 * time.Second), nil
	}

	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("代理地址解析失败: %w", err)
	}

	transport := &http.Transport{}

	switch u.Scheme {
	case "http", "https":
		transport.Proxy = http.ProxyURL(u)
	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", u.Host, nil, proxy.Direct)
		if err != nil {
			return nil, nil, fmt.Errorf("SOCKS5代理配置失败: %w", err)
		}
		transport.DialContext = dialer.(proxy.ContextDialer).DialContext
	default:
		return nil, nil, fmt.Errorf("不支持的代理协议: %s（支持 http / socks5）", u.Scheme)
	}

	fmt.Println("使用代理:", cfg.URL)
	return &http.Client{Transport: transport, Timeout: 30 * time.Second}, newClient(30 * time.Second), nil
}
