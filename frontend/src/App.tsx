import { useState } from 'react'
import FlowEditor from './FlowEditor'
import './App.css'

type Tab = 'flow' | 'agents' | 'channels' | 'settings'

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('flow')

  return (
    <div className="app">
      <header className="header">
        <h1>🚀 AgentFlow</h1>
        <nav>
          <button 
            className={activeTab === 'flow' ? 'active' : ''} 
            onClick={() => setActiveTab('flow')}
          >
            流程编排
          </button>
          <button 
            className={activeTab === 'agents' ? 'active' : ''} 
            onClick={() => setActiveTab('agents')}
          >
            智能体
          </button>
          <button 
            className={activeTab === 'channels' ? 'active' : ''} 
            onClick={() => setActiveTab('channels')}
          >
            渠道
          </button>
          <button 
            className={activeTab === 'settings' ? 'active' : ''} 
            onClick={() => setActiveTab('settings')}
          >
            设置
          </button>
        </nav>
      </header>
      
      <main className="main">
        {activeTab === 'flow' && <FlowEditor />}
        {activeTab === 'agents' && <AgentsPanel />}
        {activeTab === 'channels' && <ChannelsPanel />}
        {activeTab === 'settings' && <SettingsPanel />}
      </main>
    </div>
  )
}

function AgentsPanel() {
  return (
    <div className="panel">
      <h2>智能体管理</h2>
      <p>创建和管理AI智能体</p>
      {/* TODO: Agent CRUD */}
    </div>
  )
}

function ChannelsPanel() {
  return (
    <div className="panel">
      <h2>渠道管理</h2>
      <p>配置消息接收渠道</p>
      {/* TODO: Channel CRUD */}
    </div>
  )
}

function SettingsPanel() {
  return (
    <div className="panel">
      <h2>系统设置</h2>
      <p>配置API密钥、大模型参数等</p>
      {/* TODO: Settings */}
    </div>
  )
}

export default App
