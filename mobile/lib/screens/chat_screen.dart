import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/chat_provider.dart';

class ChatScreen extends StatefulWidget {
  const ChatScreen({super.key});

  @override
  State<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
  final _inputController = TextEditingController();
  final _scrollController = ScrollController();
  String _userId = 'user-001';

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<ChatProvider>().loadConversations(_userId);
    });
  }

  @override
  void dispose() {
    _inputController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void _scrollToBottom() {
    if (_scrollController.hasClients) {
      _scrollController.animateTo(
        _scrollController.position.maxScrollExtent,
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeOut,
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Chat'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () {
              context.read<ChatProvider>().createConversation(_userId, null);
            },
          ),
        ],
      ),
      body: Consumer<ChatProvider>(
        builder: (context, chat, _) {
          if (chat.currentConversation == null) {
            return _buildConversationList(chat);
          }
          return _buildChatView(chat);
        },
      ),
    );
  }

  Widget _buildConversationList(ChatProvider chat) {
    if (chat.conversations.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.chat_bubble_outline, size: 64, color: Colors.grey),
            const SizedBox(height: 16),
            const Text('No conversations yet'),
            const SizedBox(height: 8),
            ElevatedButton.icon(
              onPressed: () => chat.createConversation(_userId, null),
              icon: const Icon(Icons.add),
              label: const Text('New Chat'),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      itemCount: chat.conversations.length,
      itemBuilder: (context, index) {
        final conv = chat.conversations[index];
        return ListTile(
          leading: const CircleAvatar(child: Icon(Icons.chat)),
          title: Text(conv.title),
          subtitle: Text(conv.lastMessage, maxLines: 1, overflow: TextOverflow.ellipsis),
          trailing: Text(
            _formatTime(conv.updatedAt),
            style: const TextStyle(fontSize: 12, color: Colors.grey),
          ),
          onTap: () => chat.selectConversation(conv),
        );
      },
    );
  }

  Widget _buildChatView(ChatProvider chat) {
    WidgetsBinding.instance.addPostFrameCallback((_) => _scrollToBottom());

    return Column(
      children: [
        // Messages
        Expanded(
          child: chat.messages.isEmpty
              ? const Center(child: Text('Start a conversation...'))
              : ListView.builder(
                  controller: _scrollController,
                  padding: const EdgeInsets.all(16),
                  itemCount: chat.messages.length,
                  itemBuilder: (context, index) {
                    final msg = chat.messages[index];
                    final isUser = msg.sender == 'user';
                    return Align(
                      alignment: isUser ? Alignment.centerRight : Alignment.centerLeft,
                      child: Container(
                        margin: const EdgeInsets.only(bottom: 8),
                        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                        decoration: BoxDecoration(
                          color: isUser ? Theme.of(context).primaryColor : Colors.grey[200],
                          borderRadius: BorderRadius.circular(20),
                        ),
                        child: Text(
                          msg.content,
                          style: TextStyle(color: isUser ? Colors.white : Colors.black),
                        ),
                      ),
                    );
                  },
                ),
        ),
        
        // Input
        Container(
          padding: const EdgeInsets.all(8),
          decoration: BoxDecoration(
            color: Theme.of(context).cardColor,
            boxShadow: [
              BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 10),
            ],
          ),
          child: Row(
            children: [
              Expanded(
                child: TextField(
                  controller: _inputController,
                  decoration: const InputDecoration(
                    hintText: 'Type a message...',
                    border: InputBorder.none,
                    contentPadding: EdgeInsets.symmetric(horizontal: 16),
                  ),
                  maxLines: null,
                  textInputAction: TextInputAction.send,
                  onSubmitted: (_) => _sendMessage(chat),
                ),
              ),
              IconButton(
                icon: const Icon(Icons.send),
                onPressed: () => _sendMessage(chat),
              ),
            ],
          ),
        ),
      ],
    );
  }

  void _sendMessage(ChatProvider chat) {
    final text = _inputController.text.trim();
    if (text.isEmpty || chat.currentConversation == null) return;
    
    chat.sendMessage(chat.currentConversation!.id, text);
    _inputController.clear();
  }

  String _formatTime(DateTime time) {
    final now = DateTime.now();
    final diff = now.difference(time);
    if (diff.inMinutes < 60) return '${diff.inMinutes}m';
    if (diff.inHours < 24) return '${diff.inHours}h';
    return '${diff.inDays}d';
  }
}
