import 'package:flutter/material.dart';

class AgentsScreen extends StatefulWidget {
  const AgentsScreen({super.key});

  @override
  State<AgentsScreen> createState() => _AgentsScreenState();
}

class _AgentsScreenState extends State<AgentsScreen> {
  final _agents = [
    {'id': 1, 'name': 'Assistant', 'model': 'GPT-4', 'status': 'active'},
    {'id': 2, 'name': 'Translator', 'model': 'Claude-3', 'status': 'active'},
    {'id': 3, 'name': 'Coder', 'model': 'GPT-4', 'status': 'inactive'},
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Agents'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => _showCreateDialog(context),
          ),
        ],
      ),
      body: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: _agents.length,
        itemBuilder: (context, index) {
          final agent = _agents[index];
          return Card(
            margin: const EdgeInsets.only(bottom: 12),
            child: ListTile(
              leading: CircleAvatar(
                backgroundColor: agent['status'] == 'active' 
                    ? Colors.green 
                    : Colors.grey,
                child: const Icon(Icons.smart_toy, color: Colors.white),
              ),
              title: Text(agent['name']!),
              subtitle: Text('Model: ${agent['model']}'),
              trailing: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Switch(
                    value: agent['status'] == 'active',
                    onChanged: (value) {},
                  ),
                  IconButton(
                    icon: const Icon(Icons.edit_outlined),
                    onPressed: () {},
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  void _showCreateDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Create New Agent'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              decoration: const InputDecoration(
                labelText: 'Agent Name',
                hintText: 'Enter agent name',
              ),
            ),
            const SizedBox(height: 16),
            DropdownButtonFormField<String>(
              decoration: const InputDecoration(labelText: 'Model'),
              items: const [
                DropdownMenuItem(value: 'gpt-4', child: Text('GPT-4')),
                DropdownMenuItem(value: 'claude-3', child: Text('Claude 3')),
                DropdownMenuItem(value: 'glm-4', child: Text('GLM-4')),
                DropdownMenuItem(value: 'kimi', child: Text('Kimi')),
              ],
              onChanged: (value) {},
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Create'),
          ),
        ],
      ),
    );
  }
}
