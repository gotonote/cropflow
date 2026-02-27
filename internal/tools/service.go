package tools

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// ToolType 工具类型
type ToolType string

const (
	ToolBrowser    ToolType = "browser"
	ToolSearch     ToolType = "search"
	ToolFetch      ToolType = "fetch"
	ToolCalculator ToolType = "calculator"
	ToolCode       ToolType = "code"
	ToolTime       ToolType = "time"
	ToolWeather    ToolType = "weather"
	ToolWiki       ToolType = "wiki"
)

// Tool 工具定义
type Tool struct {
	Name        string      `json:"name"`
	Type        ToolType    `json:"type"`
	Description string      `json:"description"`
	Params      []Param     `json:"params"`
}

// Param 参数定义
type Param struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Default     string `json:"default"`
}

// Result 工具执行结果
type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Service 工具服务
type Service struct {
	tools    map[string]*Tool
	browsers map[string]*BrowserContext // 浏览器上下文
}

// BrowserContext 浏览器上下文
type BrowserContext struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

// NewService 创建工具服务
func NewService() *Service {
	s := &Service{
		tools:    make(map[string]*Tool),
		browsers: make(map[string]*BrowserContext),
	}
	s.registerTools()
	return s
}

// registerTools 注册所有工具
func (s *Service) registerTools() {
	// 浏览器工具
	s.tools["browser_navigate"] = &Tool{
		Name:        "browser_navigate",
		Type:        ToolBrowser,
		Description: "导航到指定URL",
		Params: []Param{
			{Name: "url", Type: "string", Required: true, Description: "目标URL"},
			{Name: "wait", Type: "number", Required: false, Description: "等待秒数"},
		},
	}

	s.tools["browser_click"] = &Tool{
		Name:        "browser_click",
		Type:        ToolBrowser,
		Description: "点击页面元素",
		Params: []Param{
			{Name: "selector", Type: "string", Required: true, Description: "CSS选择器"},
		},
	}

	s.tools["browser_input"] = &Tool{
		Name:        "browser_input",
		Type:        ToolBrowser,
		Description: "输入文本到输入框",
		Params: []Param{
			{Name: "selector", Type: "string", Required: true, Description: "CSS选择器"},
			{Name: "text", Type: "string", Required: true, Description: "输入文本"},
		},
	}

	s.tools["browser_screenshot"] = &Tool{
		Name:        "browser_screenshot",
		Type:        ToolBrowser,
		Description: "截图",
		Params: []Param{
			{Name: "fullPage", Type: "boolean", Required: false, Description: "是否截取整页"},
		},
	}

	s.tools["browser_extract"] = &Tool{
		Name:        "browser_extract",
		Type:        ToolBrowser,
		Description: "提取页面内容",
		Params: []Param{
			{Name: "selector", Type: "string", Required: false, Description: "CSS选择器"},
			{Name: "pattern", Type: "string", Required: false, Description: "正则表达式"},
		},
	}

	// 搜索工具
	s.tools["web_search"] = &Tool{
		Name:        "web_search",
		Type:        ToolSearch,
		Description: "搜索网页",
		Params: []Param{
			{Name: "query", Type: "string", Required: true, Description: "搜索关键词"},
			{Name: "limit", Type: "number", Required: false, Default: "5", Description: "结果数量"},
		},
	}

	// 获取网页
	s.tools["fetch_url"] = &Tool{
		Name:        "fetch_url",
		Type:        ToolFetch,
		Description: "获取网页内容",
		Params: []Param{
			{Name: "url", Type: "string", Required: true, Description: "目标URL"},
			{Name: "maxChars", Type: "number", Required: false, Default: "5000", Description: "最大字符数"},
		},
	}

	// 计算器
	s.tools["calculator"] = &Tool{
		Name:        "calculator",
		Type:        ToolCalculator,
		Description: "数学计算",
		Params: []Param{
			{Name: "expression", Type: "string", Required: true, Description: "数学表达式"},
		},
	}

	// 代码执行
	s.tools["run_code"] = &Tool{
		Name:        "run_code",
		Type:        ToolCode,
		Description: "执行代码",
		Params: []Param{
			{Name: "language", Type: "string", Required: true, Description: "语言: python/javascript/bash"},
			{Name: "code", Type: "string", Required: true, Description: "代码内容"},
			{Name: "timeout", Type: "number", Required: false, Default: "30", Description: "超时秒数"},
		},
	}

	// 时间
	s.tools["get_time"] = &Tool{
		Name:        "get_time",
		Type:        ToolTime,
		Description: "获取当前时间",
		Params: []Param{
			{Name: "timezone", Type: "string", Required: false, Description: "时区"},
		},
	}

	// 天气
	s.tools["get_weather"] = &Tool{
		Name:        "get_weather",
		Type:        ToolWeather,
		Description: "获取天气",
		Params: []Param{
			{Name: "city", Type: "string", Required: true, Description: "城市名称"},
		},
	}

	// Wiki搜索
	s.tools["wiki_search"] = &Tool{
		Name:        "wiki_search",
		Type:        ToolWiki,
		Description: "维基百科搜索",
		Params: []Param{
			{Name: "query", Type: "string", Required: true, Description: "搜索词"},
		},
	}
}

// Execute 执行工具
func (s *Service) Execute(ctx context.Context, toolName string, params map[string]interface{}) Result {
	tool, ok := s.tools[toolName]
	if !ok {
		return Result{Success: false, Error: fmt.Sprintf("tool not found: %s", toolName)}
	}

	var err error
	var data interface{}

	switch tool.Type {
	case ToolBrowser:
		data, err = s.execBrowser(ctx, toolName, params)
	case ToolSearch:
		data, err = s.execSearch(ctx, params)
	case ToolFetch:
		data, err = s.execFetch(ctx, params)
	case ToolCalculator:
		data, err = s.execCalculator(params)
	case ToolCode:
		data, err = s.execCode(ctx, params)
	case ToolTime:
		data, err = s.execTime(params)
	case ToolWeather:
		data, err = s.execWeather(params)
	case ToolWiki:
		data, err = s.execWiki(params)
	default:
		err = fmt.Errorf("unknown tool type: %s", tool.Type)
	}

	if err != nil {
		return Result{Success: false, Error: err.Error()}
	}

	return Result{Success: true, Data: data}
}

// ListTools 列出所有工具
func (s *Service) ListTools() []*Tool {
	tools := make([]*Tool, 0, len(s.tools))
	for _, t := range s.tools {
		tools = append(tools, t)
	}
	return tools
}

// ========== 浏览器工具实现 ==========

func (s *Service) execBrowser(ctx context.Context, toolName string, params map[string]interface{}) (interface{}, error) {
	userID := ctx.Value("user_id")
	if userID == nil {
		userID = "default"
	}
	userIDStr := userID.(string)

	switch toolName {
	case "browser_navigate":
		url, _ := params["url"].(string)
		if url == "" {
			return nil, fmt.Errorf("url is required")
		}
		return s.browserNavigate(userIDStr, url)

	case "browser_click":
		selector, _ := params["selector"].(string)
		if selector == "" {
			return nil, fmt.Errorf("selector is required")
		}
		return s.browserClick(userIDStr, selector)

	case "browser_input":
		selector, _ := params["selector"].(string)
		text, _ := params["text"].(string)
		if selector == "" || text == "" {
			return nil, fmt.Errorf("selector and text are required")
		}
		return s.browserInput(userIDStr, selector, text)

	case "browser_screenshot":
		fullPage, _ := params["fullPage"].(bool)
		return s.browserScreenshot(userIDStr, fullPage)

	case "browser_extract":
		selector, _ := params["selector"].(string)
		pattern, _ := params["pattern"].(string)
		return s.browserExtract(userIDStr, selector, pattern)

	default:
		return nil, fmt.Errorf("unknown browser tool: %s", toolName)
	}
}

// browserNavigate 导航到URL
func (s *Service) browserNavigate(userID, url string) (interface{}, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ctx, _ = chromedp.NewContext(ctx, chromedp.WithExecAllocator(opts))

	// 保存浏览器上下文
	s.browsers[userID] = &BrowserContext{Ctx: ctx, Cancel: cancel}

	var title string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Title(&title),
	)
	if err != nil {
		return nil, err
	}

	return map[string]string{"title": title, "url": url}, nil
}

// browserClick 点击元素
func (s *Service) browserClick(userID, selector string) (interface{}, error) {
	browserCtx, ok := s.browsers[userID]
	if !ok {
		return nil, fmt.Errorf("no browser session for user: %s", userID)
	}

	err := chromedp.Run(browserCtx.Ctx,
		chromedp.Click(selector),
	)
	if err != nil {
		return nil, err
	}

	return map[string]string{"action": "click", "selector": selector}, nil
}

// browserInput 输入文本
func (s *Service) browserInput(userID, selector, text string) (interface{}, error) {
	browserCtx, ok := s.browsers[userID]
	if !ok {
		return nil, fmt.Errorf("no browser session for user: %s", userID)
	}

	err := chromedp.Run(browserCtx.Ctx,
		chromedp.SetValue(selector, text),
	)
	if err != nil {
		return nil, err
	}

	return map[string]string{"action": "input", "selector": selector, "text": text}, nil
}

// browserScreenshot 截图
func (s *Service) browserScreenshot(userID string, fullPage bool) (interface{}, error) {
	browserCtx, ok := s.browsers[userID]
	if !ok {
		return nil, fmt.Errorf("no browser session for user: %s", userID)
	}

	var buf []byte
	err := chromedp.Run(browserCtx.Ctx,
		chromedp.FullScreenshot(&buf, fullPage),
	)
	if err != nil {
		return nil, err
	}

	// 返回Base64编码的图片
	base64Str := base64.StdEncoding.EncodeToString(buf)
	return map[string]string{"image": base64Str}, nil
}

// browserExtract 提取页面内容
func (s *Service) browserExtract(userID, selector, pattern string) (interface{}, error) {
	browserCtx, ok := s.browsers[userID]
	if !ok {
		return nil, fmt.Errorf("no browser session for user: %s", userID)
	}

	var result string
	if selector != "" {
		err := chromedp.Run(browserCtx.Ctx,
			chromedp.Text(selector, &result),
		)
		if err != nil {
			return nil, err
		}
	} else {
		err := chromedp.Run(browserCtx.Ctx,
			chromedp.Body(&result),
		)
		if err != nil {
			return nil, err
		}
	}

	// 如果有正则表达式，进行匹配
	if pattern != "" {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(result)
		if len(matches) > 0 {
			result = matches[0]
		}
	}

	return map[string]string{"content": result}, nil
}

// CloseBrowser 关闭浏览器
func (s *Service) CloseBrowser(userID string) {
	if browserCtx, ok := s.browsers[userID]; ok {
		browserCtx.Cancel()
		delete(s.browsers, userID)
	}
}

// ========== 搜索工具 ==========

func (s *Service) execSearch(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	query, _ := params["query"].(string)
	limit, _ := params["limit"].(float64)
	if limit == 0 {
		limit = 5
	}

	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	// TODO: 集成真实搜索API (Brave/Serper)
	// 这里返回模拟数据
	results := []map[string]string{
		{"title": "搜索结果 1", "url": "https://example.com/1", "snippet": "这是搜索结果1的摘要"},
		{"title": "搜索结果 2", "url": "https://example.com/2", "snippet": "这是搜索结果2的摘要"},
	}

	// 限制数量
	if int(limit) < len(results) {
		results = results[:int(limit)]
	}

	return map[string]interface{}{"results": results, "query": query}, nil
}

// ========== 获取网页 ==========

func (s *Service) execFetch(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	url, _ := params["url"].(string)
	maxChars, _ := params["maxChars"].(float64)
	if maxChars == 0 {
		maxChars = 5000
	}

	if url == "" {
		return nil, fmt.Errorf("url is required")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	content := string(body)
	if len(content) > int(maxChars) {
		content = content[:int(maxChars)] + "..."
	}

	return map[string]interface{}{
		"url":     url,
		"status":  resp.StatusCode,
		"content": content,
	}, nil
}

// ========== 计算器 ==========

func (s *Service) execCalculator(params map[string]interface{}) (interface{}, error) {
	expr, _ := params["expression"].(string)
	if expr == "" {
		return nil, fmt.Errorf("expression is required")
	}

	// 简单计算实现 (生产环境使用govaluate库)
	// 清理表达式
	expr = strings.ReplaceAll(expr, " ", "")

	// 验证只包含安全字符
	if !regexp.MustCompile(`^[\d+\-*/().]+$`).MatchString(expr) {
		return nil, fmt.Errorf("invalid expression")
	}

	// 使用calc命令或简单计算
	cmd := exec.Command("bc", "-l")
	cmd.Stdin = strings.NewReader(expr)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// 备用：尝试Python
		cmd = exec.Command("python3", "-c", fmt.Sprintf("print(%s)", expr))
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("calculation failed")
		}
	}

	result := strings.TrimSpace(out.String())
	return map[string]string{"expression": expr, "result": result}, nil
}

// ========== 代码执行 ==========

func (s *Service) execCode(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	language, _ := params["language"].(string)
	code, _ := params["code"].(string)
	timeout, _ := params["timeout"].(float64)
	if timeout == 0 {
		timeout = 30
	}

	if language == "" || code == "" {
		return nil, fmt.Errorf("language and code are required")
	}

	var cmd *exec.Cmd
	var stdin io.WriteCloser
	var stdout, stderr bytes.Buffer

	switch language {
	case "python":
		cmd = exec.Command("python3", "-c", code)
	case "javascript", "js":
		cmd = exec.Command("node", "-e", code)
	case "bash", "sh":
		cmd = exec.Command("bash", "-c", code)
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	cmd.Stdin = nil
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		cmd.Process.Kill()
		return nil, fmt.Errorf("execution timeout")
	case err := <-done:
		if err != nil {
			return map[string]string{
				"stdout": stdout.String(),
				"stderr": stderr.String(),
				"error":  err.Error(),
			}, nil
		}
	}

	return map[string]string{
		"stdout": stdout.String(),
		"stderr": stderr.String(),
	}, nil
}

// ========== 时间工具 ==========

func (s *Service) execTime(params map[string]interface{}) (interface{}, error) {
	timezone, _ := params["timezone"].(string)
	
	loc := time.Local
	if timezone != "" {
		var err error
		loc, err = time.LoadLocation(timezone)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone: %s", timezone)
		}
	 time.Now().}

	now :=In(loc)
	return map[string]string{
		"time":      now.Format("2006-01-02 15:04:05"),
		"timestamp": fmt.Sprintf("%d", now.Unix()),
		"timezone":  timezone,
	}, nil
}

// ========== 天气工具 ==========

func (s *Service) execWeather(params map[string]interface{}) (interface{}, error) {
	city, _ := params["city"].(string)
	if city == "" {
		return nil, fmt.Errorf("city is required")
	}

	// TODO: 集成真实天气API
	return map[string]string{
		"city":      city,
		"weather":   "晴",
		"temp":      "25°C",
		"humidity":  "60%",
		"wind":      "东北风3级",
	}, nil
}

// ========== Wiki工具 ==========

func (s *Service) execWiki(params map[string]interface{}) (interface{}, error) {
	query, _ := params["query"].(string)
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	// 模拟Wiki搜索
	return map[string]string{
		"query": query,
		"result": fmt.Sprintf("关于%s的维基百科摘要...", query),
		"url":    fmt.Sprintf("https://zh.wikipedia.org/wiki/%s", query),
	}, nil
}

// ========== OpenAI工具格式转换 ==========

// ToOpenAIFormat 转换为OpenAI工具格式
func (s *Service) ToOpenAIFormat() []map[string]interface{} {
	tools := make([]map[string]interface{}, 0)
	for _, t := range s.tools {
		props := make(map[string]interface{})
		required := []string{}
		
		for _, p := range t.Params {
			props[p.Name] = map[string]string{
				"type":        p.Type,
				"description": p.Description,
			}
			if p.Default != "" {
				props[p.Name].(map[string]string)["default"] = p.Default
			}
			if p.Required {
				required = append(required, p.Name)
			}
		}

		tools = append(tools, map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        t.Name,
				"description": t.Description,
				"parameters": map[string]interface{}{
					"type":       "object",
					"properties": props,
					"required":   required,
				},
			},
		})
	}
	return tools
}

// InitBrowser 初始化浏览器 (延迟加载chromedp)
func init() {
	// 检查chromedp是否可用
	if _, err := exec.LookPath("chromedp"); err != nil {
		fmt.Println("Warning: chromedp not found, browser tools will not work")
	}
}

// 从环境变量获取API Keys
func GetAPIKeys() (braveKey, serperKey, weatherKey string) {
	braveKey = os.Getenv("BRAVE_API_KEY")
	serperKey = os.Getenv("SERPER_API_KEY")
	weatherKey = os.Getenv("WEATHER_API_KEY")
	return
}
