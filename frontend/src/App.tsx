import { useState } from 'react'
import FlowEditor from './FlowEditor'
import './App.css'

type Tab = 'home' | 'chat' | 'flow' | 'agents' | 'tools' | 'logs' | 'settings'

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('home')

  return (
    <div className="app">
      <header className="header">
        <h1>🚀 CorpFlow</h1>
        <nav>
          <button className={activeTab === 'home' ? 'active' : ''} onClick={() => setActiveTab('home')}>🏠 Home</button>
          <button className={activeTab === 'chat' ? 'active' : ''} onClick={() => setActiveTab('chat')}>💬 Chat</button>
          <button className={activeTab === 'flow' ? 'active' : ''} onClick={() => setActiveTab('flow')}>🔀 Flows</button>
          <button className={activeTab === 'agents' ? 'active' : ''} onClick={() => setActiveTab('agents')}>🤖 Agents</button>
          <button className={activeTab === 'tools' ? 'active' : ''} onClick={() => setActiveTab('tools')}>🔧 Tools</button>
          <button className={activeTab === 'logs' ? 'active' : ''} onClick={() => setActiveTab('logs')}>📋 Logs</button>
          <button className={activeTab === 'settings' ? 'active' : ''} onClick={() => setActiveTab('settings')}>⚙️</button>
        </nav>
      </header>
      <main className="main">
        {activeTab === 'home' && <HomePanel />}
        {activeTab === 'chat' && <ChatPanel />}
        {activeTab === 'flow' && <FlowPanel />}
        {activeTab === 'agents' && <AgentsPanel />}
        {activeTab === 'tools' && <ToolsPanel />}
        {activeTab === 'logs' && <LogsPanel />}
        {activeTab === 'settings' && <SettingsPanel />}
      </main>
    </div>
  )
}

// Quick Templates
const templates = [
  { id: 'simple-chat', name: '💬 Simple Chat', desc: 'Basic AI conversation' },
  { id: 'multi-voting', name: '🗳️ Multi-Model Vote', desc: 'Multiple AI vote' },
  { id: 'research', name: '🔍 Research', desc: 'Search & analyze' },
  { id: 'customer-service', name: '🎧 Customer Service', desc: 'AI support bot' },
  { id: 'code-review', name: '📝 Code Review', desc: 'Automated review' },
  { id: 'content', name: '✍️ Content', desc: 'Social media' },
]

function HomePanel() {
  return (
    <div className="home-container">
      <section className="welcome-section">
        <h2>Welcome to CorpFlow</h2>
        <p>Multi-Agent Collaboration Platform</p>
        <div className="quick-actions">
          <button className="btn-primary">🚀 Start Chat</button>
          <button className="btn-secondary">➕ Create Flow</button>
        </div>
      </section>

      <section>
        <h3>⚡ Quick Templates</h3>
        <div className="templates-grid">
          {templates.map(t => (
            <div key={t.id} className="template-card">
              <div className="template-icon">🔗</div>
              <div className="template-name">{t.name}</div>
              <div className="template-desc">{t.desc}</div>
            </div>
          ))}
        </div>
      </section>

      <section>
        <h3>✨ Features</h3>
        <div className="features-grid">
          <div className="feature-card"><span>🤖</span><h4>AI Agents</h4><p>GPT-4, Claude, GLM-4, Kimi, Qwen, DeepSeek</p></div>
          <div className="feature-card"><span>🗳️</span><h4>Voting</h4><p>Multi-model voting & consensus</p></div>
          <div className="feature-card"><span>🔀</span><h4>Flows</h4><p>Visual workflow automation</p></div>
          <div className="feature-card"><span>💬</span><h4>Channels</h4><p>Feishu, WeChat, Telegram, Discord</p></div>
          <div className="feature-card"><span>🔧</span><h4>Tools</h4><p>Shell, Git, Code Review, Test Gen</p></div>
          <div className="feature-card"><span>📋</span><h4>Logs</h4><p>Execution tracking & debugging</p></div>
        </div>
      </section>
    </div>
  )
}

function ChatPanel() {
  const [messages, setMessages] = useState<{role: string, content: string}[]>([])
  const [input, setInput] = useState('')

  const send = () => {
    if (!input.trim()) return
    setMessages([...messages, {role: 'user', content: input}])
    setInput('')
    setTimeout(() => {
      setMessages(prev => [...prev, {role: 'bot', content: 'Configure API key in Settings to start!'}])
    }, 500)
  }

  return (
    <div className="chat-container">
      <div className="chat-messages">
        {messages.length === 0 ? (
          <div className="chat-empty"><p>💬 Start conversation</p><p className="tip">Configure API key first</p></div>
        ) : messages.map((m, i) => (
          <div key={i} className={`message ${m.role}`}>
            <span className="msg-role">{m.role === 'user' ? '👤' : '🤖'}</span>
            <span className="msg-content">{m.content}</span>
          </div>
        ))}
      </div>
      <div className="chat-input">
        <input value={input} onChange={e => setInput(e.target.value)} onKeyPress={e => e.key === 'Enter' && send()} placeholder="Type message..." />
        <button onClick={send}>Send</button>
      </div>
    </div>
  )
}

function FlowPanel() {
  return <FlowEditor />
}

function AgentsPanel() {
  return (
    <div className="agents-container">
      <h3>🤖 AI Agents</h3>
      <div className="agents-list">
        <div className="agent-item"><span>🤖</span><div><h4>Assistant</h4><p>Model: GLM-4</p></div></div>
      </div>
      <button className="btn-primary">+ Add Agent</button>
    </div>
  )
}

// ============ Tools Panel - New! ============
const builtInTools = [
  { id: 'shell', name: 'Shell', desc: 'Execute shell commands', icon: '💻' },
  { id: 'git', name: 'Git', desc: 'Git operations (commit/push)', icon: '📦' },
  { id: 'web_search', name: 'Web Search', desc: 'Search the web', icon: '🔍' },
  { id: 'web_fetch', name: 'Web Fetch', desc: 'Get web page content', icon: '🌐' },
  { id: 'file_read', name: 'File Read', desc: 'Read file content', icon: '📄' },
  { id: 'file_write', name: 'File Write', desc: 'Write file content', icon: '✏️' },
  { id: 'code_review', name: 'Code Review', desc: 'AI code review', icon: '📝' },
  { id: 'test_gen', name: 'Test Gen', desc: 'Generate unit tests', icon: '🧪' },
  { id: 'calculator', name: 'Calculator', desc: 'Math calculations', icon: '🧮' },
]

function ToolsPanel() {
  const [selectedTool, setSelectedTool] = useState<string | null>(null)
  const [toolInput, setToolInput] = useState('')
  const [toolOutput, setToolOutput] = useState('')

  const runTool = () => {
    if (!selectedTool) return
    setToolOutput(`Running ${selectedTool}...\n\n[Configure backend to execute]`)
  }

  return (
    <div className="tools-container">
      <h3>🔧 Tool Marketplace</h3>
      <p className="tip">Built-in tools for your agents</p>
      
      <div className="tools-grid">
        {builtInTools.map(tool => (
          <div key={tool.id} className={`tool-card ${selectedTool === tool.id ? 'selected' : ''}`} onClick={() => setSelectedTool(tool.id)}>
            <span className="tool-icon">{tool.icon}</span>
            <div className="tool-name">{tool.name}</div>
            <div className="tool-desc">{tool.desc}</div>
          </div>
        ))}
      </div>

      {selectedTool && (
        <div className="tool-runner">
          <h4>Run: {selectedTool}</h4>
          <textarea value={toolInput} onChange={e => setToolInput(e.target.value)} placeholder='{"key": "value"}' rows={4} />
          <button className="btn-primary" onClick={runTool}>▶️ Run</button>
          {toolOutput && <pre className="tool-output">{toolOutput}</pre>}
        </div>
      )}
    </div>
  )
}

// ============ Logs Panel - New! ============
const mockLogs = [
  { id: '1', flow: 'Simple Chat', status: 'success', time: '2 min ago', duration: '1.2s' },
  { id: '2', flow: 'Code Review', status: 'success', time: '5 min ago', duration: '3.5s' },
  { id: '3', flow: 'Research', status: 'failed', time: '10 min ago', duration: '0.5s', error: 'API key missing' },
  { id: '4', flow: 'Multi-Model Vote', status: 'running', time: 'now', duration: '-' },
]

function LogsPanel() {
  const [selectedLog, setSelectedLog] = useState<string | null>(null)

  return (
    <div className="logs-container">
      <h3>📋 Execution Logs</h3>
      <div className="logs-stats">
        <span className="stat">📊 Total: 42</span>
        <span className="stat success">✅ Success: 38</span>
        <span className="stat error">❌ Failed: 3</span>
        <span className="stat">⏱️ Avg: 2.1s</span>
      </div>

      <div className="logs-list">
        {mockLogs.map(log => (
          <div key={log.id} className={`log-item ${selectedLog === log.id ? 'selected' : ''}`} onClick={() => setSelectedLog(log.id)}>
            <span className={`status-dot ${log.status}`}></span>
            <div className="log-info">
              <div className="log-flow">{log.flow}</div>
              <div className="log-meta">{log.time} · {log.duration}</div>
            </div>
            {log.error && <span className="log-error">{log.error}</span>}
          </div>
        ))}
      </div>

      {selectedLog && (
        <div className="log-detail">
          <h4>Log Detail</h4>
          <pre>{`ID: ${selectedLog}
Flow: Code Review
Status: Success
Duration: 3.5s

Steps:
✅ Trigger (message) - 0.1s
✅ Agent (GPT-4) - 2.1s  
✅ Code Review - 1.2s
✅ Output - 0.1s`}</pre>
        </div>
      )}
    </div>
  )
}

function SettingsPanel() {
  return (
    <div className="settings-container">
      <h3>⚙️ Settings</h3>
      <div className="settings-section">
        <h4>🔑 API Keys</h4>
        <div className="api-key-input"><label>OpenAI (GPT-4)</label><input type="password" placeholder="sk-..." /></div>
        <div className="api-key-input"><label>Zhipu (GLM-4)</label><input type="password" placeholder="API key" /></div>
        <div className="api-key-input"><label>Kimi</label><input type="password" placeholder="API key" /></div>
        <div className="api-key-input"><label>Qwen</label><input type="password" placeholder="API key" /></div>
        <div className="api-key-input"><label>DeepSeek</label><input type="password" placeholder="API key" /></div>
        <button className="btn-save">Save Keys</button>
      </div>
      <div className="settings-section">
        <h4>🎯 Default Model</h4>
        <select><option>GLM-4 (Recommended)</option><option>GPT-4</option><option>Kimi</option></select>
      </div>
      <div className="settings-section">
        <h4>🗳️ Multi-Model Voting</h4>
        <label className="toggle"><input type="checkbox" defaultChecked /><span>Enable voting</span></label>
      </div>
    </div>
  )
}

export default App
