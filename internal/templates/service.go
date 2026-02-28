// CorpFlow Templates - Ready to use workflow templates

package templates

// Template 1: Simple Chat
const SimpleChat = `{
  "name": "Simple Chat",
  "description": "Basic AI chat workflow",
  "nodes": [
    {"id": "1", "type": "trigger", "data": {"label": "用户消息", "triggerType": "message"}},
    {"id": "2", "type": "agent", "data": {"label": "AI助手", "model": "glm-4"}},
    {"id": "3", "type": "output", "data": {"label": "返回回复"}}
  ],
  "edges": [
    {"id": "e1-2", "source": "1", "target": "2"},
    {"id": "e2-3", "source": "2", "target": "3"}
  ]
}`

// Template 2: Multi-Agent Collaboration
const MultiAgentCollaboration = `{
  "name": "Multi-Agent Collaboration",
  "description": "CEO → Manager → Worker collaborative problem solving",
  "nodes": [
    {"id": "1", "type": "trigger", "data": {"label": "用户任务"}},
    {"id": "2", "type": "agent", "data": {"label": "CEO", "role": "ceo", "prompt": "分析任务，分解为子任务"}},
    {"id": "3", "type": "agent", "data": {"label": "Manager", "role": "manager", "prompt": "分配子任务给Worker"}},
    {"id": "4", "type": "agent", "data": {"label": "Worker", "role": "worker", "prompt": "执行具体任务"}},
    {"id": "5", "type": "agent", "data": {"label": "CEO", "role": "ceo_final", "prompt": "汇总结果，给出最终方案"}},
    {"id": "6", "type": "output", "data": {"label": "输出最终答案"}}
  ],
  "edges": [
    {"id": "e1-2", "source": "1", "target": "2"},
    {"id": "e2-3", "source": "2", "target": "3"},
    {"id": "e3-4", "source": "3", "target": "4"},
    {"id": "e4-5", "source": "4", "target": "5"},
    {"id": "e5-6", "source": "5", "target": "6"}
  ]
}`

// Template 3: Research Assistant
const ResearchAssistant = `{
  "name": "Research Assistant",
  "description": "Search, analyze and summarize information",
  "nodes": [
    {"id": "1", "type": "trigger", "data": {"label": "研究主题"}},
    {"id": "2", "type": "tool", "data": {"label": "搜索", "toolName": "web_search"}},
    {"id": "3", "type": "agent", "data": {"label": "分析", "model": "glm-4", "prompt": "分析搜索结果"}},
    {"id": "4", "type": "tool", "data": {"label": "获取网页", "toolName": "fetch_url"}},
    {"id": "5", "type": "agent", "data": {"label": "总结", "model": "glm-4", "prompt": "生成研究报告"}},
    {"id": "6", "type": "output", "data": {"label": "输出报告"}}
  ],
  "edges": [
    {"id": "e1-2", "source": "1", "target": "2"},
    {"id": "e2-3", "source": "2", "target": "3"},
    {"id": "e3-4", "source": "3", "target": "4"},
    {"id": "e4-5", "source": "4", "target": "5"},
    {"id": "e5-6", "source": "5", "target": "6"}
  ]
}`

// Template 4: Customer Service
const CustomerService = `{
  "name": "Customer Service",
  "description": "AI customer support with knowledge base",
  "nodes": [
    {"id": "1", "type": "trigger", "data": {"label": "客户咨询"}},
    {"id": "2", "type": "agent", "data": {"label": "理解问题", "model": "glm-4"}},
    {"id": "3", "type": "condition", "data": {"label": "是否已知问题"}},
    {"id": "4", "type": "agent", "data": {"label": "知识库回答", "model": "glm-4"}},
    {"id": "5", "type": "agent", "data": {"label": "转人工", "model": "glm-4"}},
    {"id": "6", "type": "output", "data": {"label": "回复客户"}}
  ],
  "edges": [
    {"id": "e1-2", "source": "1", "target": "2"},
    {"id": "e2-3", "source": "2", "target": "3"},
    {"id": "e3-4", "source": "3", "target": "4", "condition": "已知"},
    {"id": "e3-5", "source": "3", "target": "5", "condition": "未知"},
    {"id": "e4-6", "source": "4", "target": "6"},
    {"id": "e5-6", "source": "5", "target": "6"}
  ]
}`

// Template 5: Code Review
const CodeReview = `{
  "name": "Code Review",
  "description": "Automated code review with multiple models",
  "nodes": [
    {"id": "1", "type": "trigger", "data": {"label": "代码提交", "triggerType": "webhook"}},
    {"id": "2", "type": "tool", "data": {"label": "获取代码", "toolName": "fetch_url"}},
    {"id": "3", "type": "agent", "data": {"label": "语法检查", "model": "deepseek-coder"}},
    {"id": "4", "type": "agent", "data": {"label": "代码优化", "model": "gpt-4"}},
    {"id": "5", "type": "agent", "data": {"label": "安全审计", "model": "claude-3-sonnet"}},
    {"id": "6", "type": "condition", "data": {"label": "综合评分"}},
    {"id": "7", "type": "output", "data": {"label": "审查报告"}}
  ],
  "edges": [
    {"id": "e1-2", "source": "1", "target": "2"},
    {"id": "e2-3", "source": "2", "target": "3"},
    {"id": "e2-4", "source": "2", "target": "4"},
    {"id": "e2-5", "source": "2", "target": "5"},
    {"id": "e3-6", "source": "3", "target": "6"},
    {"id": "e4-6", "source": "4", "target": "6"},
    {"id": "e5-6", "source": "5", "target": "6"},
    {"id": "e6-7", "source": "6", "target": "7"}
  ]
}`

// Template 6: Content Creator
const ContentCreator = `{
  "name": "Content Creator",
  "description": "Generate social media content with voting",
  "nodes": [
    {"id": "1", "type": "trigger", "data": {"label": "主题输入"}},
    {"id": "2", "type": "agent", "data": {"label": "创意生成", "model": "gpt-4"}},
    {"id": "3", "type": "agent", "data": {"label": "Twitter版本", "model": "glm-4"}},
    {"id": "4", "type": "agent", "data": {"label": "微信版本", "model": "glm-4"}},
    {"id": "5", "type": "condition", "data": {"label": "选择最佳"}},
    {"id": "6", "type": "output", "data": {"label": "发布内容"}}
  ],
  "edges": [
    {"id": "e1-2", "source": "1", "target": "2"},
    {"id": "e2-3", "source": "2", "target": "3"},
    {"id": "e2-4", "source": "2", "target": "4"},
    {"id": "e3-5", "source": "3", "target": "5"},
    {"id": "e4-5", "source": "4", "target": "5"},
    {"id": "e5-6", "source": "5", "target": "6"}
  ]
}`

// Template 7: Data Analyzer
const DataAnalyzer = `{
  "name": "Data Analyzer",
  "description": "Analyze data and generate insights",
  "nodes": [
    {"id": "1", "type": "trigger", "data": {"label": "数据输入"}},
    {"id": "2", "type": "tool", "data": {"label": "数据处理", "toolName": "calculator"}},
    {"id": "3", "type": "agent", "data": {"label": "统计分析", "model": "gpt-4"}},
    {"id": "4", "type": "agent", "data": {"label": "可视化建议", "model": "gpt-4"}},
    {"id": "5", "type": "agent", "data": {"label": "生成报告", "model": "glm-4"}},
    {"id": "6", "type": "output", "data": {"label": "分析报告"}}
  ],
  "edges": [
    {"id": "e1-2", "source": "1", "target": "2"},
    {"id": "e2-3", "source": "2", "target": "3"},
    {"id": "e3-4", "source": "3", "target": "4"},
    {"id": "e4-5", "source": "4", "target": "5"},
    {"id": "e5-6", "source": "5", "target": "6"}
  ]
}`

// Template 8: News Summarizer
const NewsSummarizer = `{
  "name": "News Summarizer",
  "description": "Daily news aggregation and summary",
  "nodes": [
    {"id": "1", "type": "trigger", "data": {"label": "定时触发", "triggerType": "schedule"}},
    {"id": "2", "type": "tool", "data": {"label": "搜索科技新闻", "toolName": "web_search"}},
    {"id": "3", "type": "tool", "data": {"label": "搜索财经新闻", "toolName": "web_search"}},
    {"id": "4", "type": "agent", "data": {"label": "AI摘要", "model": "glm-4"}},
    {"id": "5", "type": "condition", "data": {"label": "是否发送"}},
    {"id": "6", "type": "output", "data": {"label": "推送给用户"}}
  ],
  "edges": [
    {"id": "e1-2", "source": "1", "target": "2"},
    {"id": "e1-3", "source": "1", "target": "3"},
    {"id": "e2-4", "source": "2", "target": "4"},
    {"id": "e3-4", "source": "3", "target": "4"},
    {"id": "e4-5", "source": "4", "target": "5"},
    {"id": "e5-6", "source": "5", "target": "6"}
  ]
}`

// Template: Superpowers-style Software Development Workflow
// Based on Superpowers (obra/superpowers) - 65k stars
// Workflow: Brainstorming → Planning → Subagent Dev → TDD → Code Review → Finish
const SuperpowersDev = `{
  "name": "Superpowers Development",
  "description": "Professional SDLC: Brainstorm → Plan → TDD → Review → Finish",
  "nodes": [
    {"id": "1", "type": "trigger", "data": {"label": "需求输入"}},
    {"id": "2", "type": "agent", "data": {"label": "Brainstorming", "role": "designer", "prompt": "通过提问澄清需求，展示设计方案"}},
    {"id": "3", "type": "agent", "data": {"label": "Planning", "role": "planner", "prompt": "创建详细实现计划，分解为2-5分钟的任务"}},
    {"id": "4", "type": "agent", "data": {"label": "Subagent Dev", "role": "developer", "prompt": "并行执行每个任务，两阶段审查"}},
    {"id": "5", "type": "agent", "data": {"label": "TDD", "role": "tester", "prompt": "RED-GREEN-REFACTOR: 先写测试，再写代码"}},
    {"id": "6", "type": "agent", "data": {"label": "Code Review", "role": "reviewer", "prompt": "按计划审查代码，报告问题"}},
    {"id": "7", "type": "agent", "data": {"label": "Finish", "role": "manager", "prompt": "验证测试，合并/PR/清理"}},
    {"id": "8", "type": "output", "data": {"label": "开发完成"}}
  ],
  "edges": [
    {"id": "e1-2", "source": "1", "target": "2"},
    {"id": "e2-3", "source": "2", "target": "3"},
    {"id": "e3-4", "source": "3", "target": "4"},
    {"id": "e4-5", "source": "4", "target": "5"},
    {"id": "e5-6", "source": "5", "target": "6"},
    {"id": "e6-7", "source": "6", "target": "7"},
    {"id": "e7-8", "source": "7", "target": "8"}
  ]
}`

// All templates
var AllTemplates = []Template{
  {"simple-chat", "Simple Chat", "Basic AI chat", SimpleChat},
  {"multi-agent-collaboration", "Multi-Agent Collaboration", "CEO + Manager + Worker collaboration", MultiAgentCollaboration},
  {"superpowers-dev", "Superpowers Dev", "Professional SDLC workflow (TDD, Code Review)", SuperpowersDev},
  {"research-assistant", "Research Assistant", "Search & analyze", ResearchAssistant},
  {"customer-service", "Customer Service", "AI support", CustomerService},
  {"code-review", "Code Review", "Automated code review", CodeReview},
  {"content-creator", "Content Creator", "Social media content", ContentCreator},
  {"data-analyzer", "Data Analyzer", "Data insights", DataAnalyzer},
  {"news-summarizer", "News Summarizer", "Daily news", NewsSummarizer},
}

type Template struct {
  ID          string `json:"id"`
  Name        string `json:"name"`
  Description string `json:"description"`
  Content     string `json:"content"`
}

func GetTemplates() []Template {
  return AllTemplates
}

func GetTemplate(id string) *Template {
  for _, t := range AllTemplates {
    if t.ID == id {
      return &t
    }
  }
  return nil
}
