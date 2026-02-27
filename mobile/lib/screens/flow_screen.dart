import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/flow_provider.dart';

class FlowScreen extends StatefulWidget {
  const FlowScreen({super.key});

  @override
  State<FlowScreen> createState() => _FlowScreenState();
}

class _FlowScreenState extends State<FlowScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<FlowProvider>().loadFlows();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Flows'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => _showCreateDialog(context),
          ),
        ],
      ),
      body: Consumer<FlowProvider>(
        builder: (context, flow, _) {
          if (flow.isLoading) {
            return const Center(child: CircularProgressIndicator());
          }
          
          if (flow.flows.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.account_tree_outlined, size: 64, color: Colors.grey),
                  const SizedBox(height: 16),
                  const Text('No flows yet'),
                  const SizedBox(height: 8),
                  ElevatedButton.icon(
                    onPressed: () => _showCreateDialog(context),
                    icon: const Icon(Icons.add),
                    label: const Text('Create Flow'),
                  ),
                ],
              ),
            );
          }
          
          return ListView.builder(
            padding: const EdgeInsets.all(16),
            itemCount: flow.flows.length,
            itemBuilder: (context, index) {
              final f = flow.flows[index];
              return Card(
                margin: const EdgeInsets.only(bottom: 12),
                child: ListTile(
                  leading: CircleAvatar(
                    backgroundColor: f.enabled ? Colors.green : Colors.grey,
                    child: const Icon(Icons.account_tree, color: Colors.white),
                  ),
                  title: Text(f.name),
                  subtitle: Text('${f.nodes.length} nodes • ${f.edges.length} edges'),
                  trailing: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      IconButton(
                        icon: const Icon(Icons.play_arrow),
                        onPressed: () => _executeFlow(context, f.id),
                      ),
                      IconButton(
                        icon: const Icon(Icons.delete_outline),
                        onPressed: () => flow.deleteFlow(f.id),
                      ),
                    ],
                  ),
                ),
              );
            },
          );
        },
      ),
    );
  }

  void _showCreateDialog(BuildContext context) {
    final controller = TextEditingController();
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Create New Flow'),
        content: TextField(
          controller: controller,
          decoration: const InputDecoration(
            labelText: 'Flow Name',
            hintText: 'Enter flow name',
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              if (controller.text.isNotEmpty) {
                context.read<FlowProvider>().createFlow(controller.text);
                Navigator.pop(context);
              }
            },
            child: const Text('Create'),
          ),
        ],
      ),
    );
  }

  void _executeFlow(BuildContext context, int id) {
    final inputController = TextEditingController();
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Execute Flow'),
        content: TextField(
          controller: inputController,
          decoration: const InputDecoration(
            labelText: 'Input',
            hintText: 'Enter input for the flow',
          ),
          maxLines: 3,
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              context.read<FlowProvider>().executeFlow(id, inputController.text);
              Navigator.pop(context);
            },
            child: const Text('Execute'),
          ),
        ],
      ),
    );
  }
}
