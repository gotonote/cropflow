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

  @override
  void initState() {
    super.initState();
    _loadSettings();
  }

  Future<void> _loadSettings() async {
    final prefs = await SharedPreferences.getInstance();
    setState(() {
      _serverController.text = prefs.getString('server_url') ?? 'http://localhost:8080';
      _darkMode = prefs.getBool('dark_mode') ?? false;
    });
  }

  Future<void> _saveSettings() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('server_url', _serverController.text);
    await prefs.setBool('dark_mode', _darkMode);
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
          _buildSectionHeader('Server / 服务器'),
          ListTile(
            leading: const Icon(Icons.dns),
            title: const Text('Server URL'),
            subtitle: Text(_serverController.text),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showServerDialog(),
          ),
          
          // Appearance
          _buildSectionHeader('Appearance / 外观'),
          SwitchListTile(
            secondary: const Icon(Icons.dark_mode),
            title: const Text('Dark Mode'),
            subtitle: const Text('深色模式'),
            value: _darkMode,
            onChanged: (value) {
              setState(() => _darkMode = value);
              _saveSettings();
            },
          ),
          
          // Model Settings
          _buildSectionHeader('AI Models / AI模型'),
          ListTile(
            leading: const Icon(Icons.model_training),
            title: const Text('Default Model'),
            subtitle: const Text('GPT-4'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {},
          ),
          ListTile(
            leading: const Icon(Icons.poll),
            title: const Text('Multi-Model Voting'),
            subtitle: const Text('Enable model voting / 启用模型投票'),
            trailing: Switch(
              value: true,
              onChanged: (value) {},
            ),
          ),
          
          // Channels
          _buildSectionHeader('Channels / 渠道'),
          ListTile(
            leading: const Icon(Icons.chat),
            title: const Text('Feishu / 飞书'),
            trailing: Switch(value: true, onChanged: (v) {}),
          ),
          ListTile(
            leading: const Icon(Icons.chat),
            title: const Text('WeChat / 微信'),
            trailing: Switch(value: true, onChanged: (v) {}),
          ),
          ListTile(
            leading: const Icon(Icons.chat),
            title: const Text('Telegram'),
            trailing: Switch(value: false, onChanged: (v) {}),
          ),
          
          // About
          _buildSectionHeader('About / 关于'),
          ListTile(
            leading: const Icon(Icons.info),
            title: const Text('Version'),
            subtitle: const Text('1.0.0'),
          ),
          ListTile(
            leading: const Icon(Icons.code),
            title: const Text('GitHub'),
            subtitle: const Text('github.com/gotonote/corpflow'),
            trailing: const Icon(Icons.open_in_new),
            onTap: () {},
          ),
          
          const SizedBox(height: 32),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: ElevatedButton(
              onPressed: _saveSettings,
              child: const Text('Save Settings'),
            ),
          ),
          const SizedBox(height: 32),
        ],
      ),
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

  @override
  void dispose() {
    _serverController.dispose();
    super.dispose();
  }
}
