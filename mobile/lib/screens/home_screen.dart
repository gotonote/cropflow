import 'package:flutter/material.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('CorpFlow'),
        centerTitle: true,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Welcome Card
            Card(
              child: Padding(
                padding: const EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Icon(Icons.smart_toy, size: 40, color: Color(0xFF667eea)),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                'Welcome to CorpFlow',
                                style: Theme.of(context).textTheme.headlineSmall,
                              ),
                              Text(
                                'Multi-Agent Collaboration Platform',
                                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                                  color: Colors.grey[600],
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 20),
            
            // Quick Actions
            Text(
              'Quick Actions',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: _QuickActionCard(
                    icon: Icons.chat,
                    title: 'New Chat',
                    subtitle: 'Start conversation',
                    onTap: () {},
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _QuickActionCard(
                    icon: Icons.account_tree,
                    title: 'New Flow',
                    subtitle: 'Create workflow',
                    onTap: () {},
                  ),
                ),
              ],
            ),
            const SizedBox(height: 20),
            
            // Statistics
            Text(
              'Statistics',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(child: _StatCard(title: 'Agents', value: '5', icon: Icons.smart_toy)),
                const SizedBox(width: 12),
                Expanded(child: _StatCard(title: 'Flows', value: '12', icon: Icons.account_tree)),
                const SizedBox(width: 12),
                Expanded(child: _StatCard(title: 'Messages', value: '1.2k', icon: Icons.chat)),
              ],
            ),
            const SizedBox(height: 20),
            
            // Recent Activity
            Text(
              'Recent Activity',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Card(
              child: ListView.separated(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                itemCount: 3,
                separatorBuilder: (_, __) => const Divider(height: 1),
                itemBuilder: (context, index) => ListTile(
                  leading: CircleAvatar(
                    backgroundColor: Color(0xFF667eea).withOpacity(0.1),
                    child: const Icon(Icons.chat, color: Color(0xFF667eea)),
                  ),
                  title: Text('Conversation ${index + 1}'),
                  subtitle: Text('Last message: 5 min ago'),
                  trailing: const Icon(Icons.chevron_right),
                ),
              ),
            ),
            const SizedBox(height: 20),
            
            // Links
            Text(
              'Links / 链接',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Card(
              child: Column(
                children: [
                  ListTile(
                    leading: const Icon(Icons.language),
                    title: const Text('Documentation'),
                    subtitle: const Text('View docs / 查看文档'),
                    trailing: const Icon(Icons.open_in_new),
                    onTap: () {},
                  ),
                  const Divider(height: 1),
                  ListTile(
                    leading: const Icon(Icons.code),
                    title: const Text('GitHub'),
                    subtitle: const Text('Source code / 源代码'),
                    trailing: const Icon(Icons.open_in_new),
                    onTap: () {},
                  ),
                  const Divider(height: 1),
                  ListTile(
                    leading: const Icon(Icons.forum),
                    title: const Text('Discord'),
                    subtitle: const Text('Community / 社区'),
                    trailing: const Icon(Icons.open_in_new),
                    onTap: () {},
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _QuickActionCard extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onTap;

  const _QuickActionCard({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            children: [
              Icon(icon, size: 32, color: Color(0xFF667eea)),
              const SizedBox(height: 8),
              Text(title, style: const TextStyle(fontWeight: FontWeight.bold)),
              Text(subtitle, style: TextStyle(fontSize: 12, color: Colors.grey[600])),
            ],
          ),
        ),
      ),
    );
  }
}

class _StatCard extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;

  const _StatCard({
    required this.title,
    required this.value,
    required this.icon,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          children: [
            Icon(icon, color: Color(0xFF667eea)),
            const SizedBox(height: 4),
            Text(value, style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
            Text(title, style: TextStyle(fontSize: 12, color: Colors.grey[600])),
          ],
        ),
      ),
    );
  }
}
