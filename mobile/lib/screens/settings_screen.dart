import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SettingsScreen extends StatefulWidget {
  const SettingsScreen({super.key});

  @override
  State<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends State<SettingsScreen> {
  final _serverController = TextEditingController(text: 'http://localhost:8080');
  bool _darkMode = false;
  bool _multiModelVoting = true;
  String _defaultModel = 'gpt-4';
  String _votingMethod = 'comprehensive';

  // Model API Keys
  final Map<String, TextEditingController> _apiKeyControllers = {};
  final Map<String, bool> _modelEnabled = {};

  final List<Map<String, String>> _models = [
    {'id': 'gpt-4', 'name': 'GPT-4', 'provider': 'OpenAI', 'keyEnv': 'OPENAI_API_KEY', 'keyName': 'openai_api_key'},
    {'id': 'gpt-3.5-turbo', 'name': 'GPT-3.5 Turbo', 'provider': 'OpenAI', 'keyEnv': 'OPENAI_API_KEY', 'keyName': 'openai_api_key'},
    {'id': 'claude-3-opus', 'name': 'Claude 3 Opus', 'provider': 'Anthropic', 'keyEnv': 'ANTHROPIC_API_KEY', 'keyName': 'anthropic_api_key'},
    {'id': 'claude-3-sonnet', 'name': 'Claude 3 Sonnet', 'provider': 'Anthropic', 'keyEnv': 'ANTHROPIC_API_KEY', 'keyName': 'anthropic_api_key'},
    {'id': 'glm-4', 'name': 'GLM-4', 'provider': 'Zhipu (智谱)', 'keyEnv': 'ZHIPU_API_KEY', 'keyName': 'zhipu_api_key'},
    {'id': 'glm-3-turbo', 'name': 'GLM-3 Turbo', 'provider': 'Zhipu (智谱)', 'keyEnv': 'ZHIPU_API_KEY', 'keyName': 'zhipu_api_key'},
    {'id': 'moonshot-v1-8k-chat', 'name': 'Kimi', 'provider': 'Moonshot (月之暗面)', 'keyEnv': 'KIMI_API_KEY', 'keyName': 'kimi_api_key'},
    {'id': 'moonshot-v1-32k-chat', 'name': 'Kimi 32K', 'provider': 'Moonshot (月之暗面)', 'keyEnv': 'KIMI_API_KEY', 'keyName': 'kimi_api_key'},
    {'id': 'qwen-turbo', 'name': 'Qwen Turbo', 'provider': 'Alibaba (通义千问)', 'keyEnv': 'DASHSCOPE_API_KEY', 'keyName': 'dashscope_api_key'},
    {'id': 'qwen-plus', 'name': 'Qwen Plus', 'provider': 'Alibaba (通义千问)', 'keyEnv': 'DASHSCOPE_API_KEY', 'keyName': 'dashscope_api_key'},
    {'id': 'qwen-max', 'name': 'Qwen Max', 'provider': 'Alibaba (通义千问)', 'keyEnv': 'DASHSCOPE_API_KEY', 'keyName': 'dashscope_api_key'},
    {'id': 'deepseek-chat', 'name': 'DeepSeek Chat', 'provider': 'DeepSeek', 'keyEnv': 'DEEPSEEK_API_KEY', 'keyName': 'deepseek_api_key'},
    {'id': 'deepseek-coder', 'name': 'DeepSeek Coder', 'provider': 'DeepSeek', 'keyEnv': 'DEEPSEEK_API_KEY', 'keyName': 'deepseek_api_key'},
    {'id': 'abab6.5s-chat', 'name': 'MiniMax', 'provider': 'MiniMax', 'keyEnv': 'MINIMAX_API_KEY', 'keyName': 'minimax_api_key'},
  ];

  @override
  void initState() {
    super.initState();
    // Initialize controllers
    for (var model in _models) {
      _apiKeyControllers[model['id']!] = TextEditingController();
      _modelEnabled[model['id']!] = false;
    }
    _loadSettings();
  }

  Future<void> _loadSettings() async {
    final prefs = await SharedPreferences.getInstance();
    setState(() {
      _serverController.text = prefs.getString('server_url') ?? 'http://localhost:8080';
      _darkMode = prefs.getBool('dark_mode') ?? false;
      _multiModelVoting = prefs.getBool('multi_model_voting') ?? true;
      _defaultModel = prefs.getString('default_model') ?? 'gpt-4';
      _votingMethod = prefs.getString('voting_method') ?? 'comprehensive';

      // Load API keys
      for (var model in _models) {
        final key = model['keyName']!;
        final apiKey = prefs.getString('key_$key') ?? '';
        _apiKeyControllers[model['id']!]!.text = apiKey;
        _modelEnabled[model['id']!] = apiKey.isNotEmpty;
      }
    });
  }

  Future<void> _saveSettings() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('server_url', _serverController.text);
    await prefs.setBool('dark_mode', _darkMode);
    await prefs.setBool('multi_model_voting', _multiModelVoting);
    await prefs.setString('default_model', _defaultModel);
    await prefs.setString('voting_method', _votingMethod);

    // Save API keys
    for (var model in _models) {
      final key = model['keyName']!;
      await prefs.setString('key_$key', _apiKeyControllers[model['id']!]!.text);
    }

    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Settings saved')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Settings'),
      ),
      body: ListView(
        children: [
          // Server Settings
          _buildSectionHeader('Server'),
          ListTile(
            leading: const Icon(Icons.dns),
            title: const Text('Server URL'),
            subtitle: Text(_serverController.text),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showServerDialog(),
          ),

          // Appearance
          _buildSectionHeader('Appearance'),
          SwitchListTile(
            secondary: const Icon(Icons.dark_mode),
            title: const Text('Dark Mode'),
            value: _darkMode,
            onChanged: (value) {
              setState(() => _darkMode = value);
              _saveSettings();
            },
          ),

          // Default Model
          _buildSectionHeader('Default AI Model'),
          ListTile(
            leading: const Icon(Icons.model_training),
            title: const Text('Default Model'),
            subtitle: Text(_defaultModel),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showDefaultModelDialog(),
          ),

          // Multi-Model Voting
          _buildSectionHeader('Multi-Model Voting'),
          SwitchListTile(
            secondary: const Icon(Icons.poll),
            title: const Text('Enable Multi-Model Voting'),
            value: _multiModelVoting,
            onChanged: (value) {
              setState(() => _multiModelVoting = value);
              _saveSettings();
            },
          ),
          if (_multiModelVoting)
            ListTile(
              leading: const Icon(Icons.analytics),
              title: const Text('Voting Method'),
              subtitle: Text(_getVotingMethodName(_votingMethod)),
              trailing: const Icon(Icons.chevron_right),
              onTap: () => _showVotingMethodDialog(),
            ),

          // API Keys Management
          _buildSectionHeader('API Keys'),
          ..._models.map((model) => _buildModelTile(model)),

          // Channels
          _buildSectionHeader('Channels'),
          SwitchListTile(
            secondary: const Icon(Icons.chat_bubble),
            title: const Text('Feishu'),
            value: true,
            onChanged: (v) {},
          ),
          SwitchListTile(
            secondary: const Icon(Icons.chat_bubble_outline),
            title: const Text('WeChat'),
            value: true,
            onChanged: (v) {},
          ),

          // About
          _buildSectionHeader('About'),
          const ListTile(
            leading: Icon(Icons.info),
            title: Text('Version'),
            subtitle: Text('1.0.0'),
          ),

          const SizedBox(height: 32),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: ElevatedButton(
              onPressed: _saveSettings,
              child: const Text('Save All Settings'),
            ),
          ),
          const SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildModelTile(Map<String, String> model) {
    return ExpansionTile(
      leading: Icon(
        _modelEnabled[model['id']!] ? Icons.check_circle : Icons.circle_outlined,
        color: _modelEnabled[model['id']!] ? Colors.green : Colors.grey,
      ),
      title: Text(model['name']!),
      subtitle: Text(model['provider']!),
      children: [
        Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              TextField(
                controller: _apiKeyControllers[model['id']],
                decoration: InputDecoration(
                  labelText: 'API Key',
                  hintText: 'Enter your ${model['provider']} API key',
                  border: const OutlineInputBorder(),
                  suffixIcon: IconButton(
                    icon: const Icon(Icons.save),
                    onPressed: () {
                      final keyName = model['keyName']!;
                      final keyValue = _apiKeyControllers[model['id']!]!.text;
                      setState(() {
                        _modelEnabled[model['id']!] = keyValue.isNotEmpty;
                      });
                      _saveSettings();
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(content: Text('${model['name']} API key saved')),
                      );
                    },
                  ),
                ),
                obscureText: true,
              ),
              const SizedBox(height: 8),
              Text(
                'Environment: ${model['keyEnv']}',
                style: TextStyle(fontSize: 12, color: Colors.grey[600]),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildSectionHeader(String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 24, 16, 8),
      child: Text(
        title,
        style: TextStyle(
          color: Theme.of(context).primaryColor,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  String _getVotingMethodName(String method) {
    switch (method) {
      case 'comprehensive':
        return 'Comprehensive';
      case 'cross':
        return 'Cross-evaluation';
      case 'length':
        return 'By Length';
      default:
        return method;
    }
  }

  void _showServerDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Server URL'),
        content: TextField(
          controller: _serverController,
          decoration: const InputDecoration(
            labelText: 'URL',
            hintText: 'http://localhost:8080',
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              _saveSettings();
              Navigator.pop(context);
            },
            child: const Text('Save'),
          ),
        ],
      ),
    );
  }

  void _showDefaultModelDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Select Default Model'),
        content: SizedBox(
          width: double.maxFinite,
          child: ListView.builder(
            shrinkWrap: true,
            itemCount: _models.length,
            itemBuilder: (context, index) {
              final model = _models[index];
              return RadioListTile<String>(
                title: Text(model['name']!),
                subtitle: Text(model['provider']!),
                value: model['id']!,
                groupValue: _defaultModel,
                onChanged: (v) {
                  setState(() => _defaultModel = v!);
                  Navigator.pop(context);
                  _saveSettings();
                },
              );
            },
          ),
        ),
      ),
    );
  }

  void _showVotingMethodDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Voting Method'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            RadioListTile<String>(
              title: const Text('Comprehensive'),
              subtitle: const Text('Accuracy/Completeness/Clarity/Creativity'),
              value: 'comprehensive',
              groupValue: _votingMethod,
              onChanged: (v) {
                setState(() => _votingMethod = v!);
                Navigator.pop(context);
                _saveSettings();
              },
            ),
            RadioListTile<String>(
              title: const Text('Cross-evaluation'),
              subtitle: const Text('Models evaluate each other'),
              value: 'cross',
              groupValue: _votingMethod,
              onChanged: (v) {
                setState(() => _votingMethod = v!);
                Navigator.pop(context);
                _saveSettings();
              },
            ),
            RadioListTile<String>(
              title: const Text('By Length'),
              subtitle: const Text('Simple length-based'),
              value: 'length',
              groupValue: _votingMethod,
              onChanged: (v) {
                setState(() => _votingMethod = v!);
                Navigator.pop(context);
                _saveSettings();
              },
            ),
          ],
        ),
      ),
    );
  }

  @override
  void dispose() {
    _serverController.dispose();
    for (var controller in _apiKeyControllers.values) {
      controller.dispose();
    }
    super.dispose();
  }
}
