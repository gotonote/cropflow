import 'package:flutter/foundation.dart';
import '../services/api_service.dart';
import '../services/websocket_service.dart';

class Message {
  final String id;
  final String type;
  final String content;
  final String sender;
  final String senderId;
  final DateTime createdAt;

  Message({
    required this.id,
    required this.type,
    required this.content,
    required this.sender,
    required this.senderId,
    required this.createdAt,
  });

  factory Message.fromJson(Map<String, dynamic> json) {
    return Message(
      id: json['id'] ?? '',
      type: json['type'] ?? 'text',
      content: json['content'] ?? '',
      sender: json['sender'] ?? 'user',
      senderId: json['sender_id'] ?? '',
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class Conversation {
  final String id;
  final String title;
  final String lastMessage;
  final DateTime updatedAt;
  List<Message> messages;

  Conversation({
    required this.id,
    required this.title,
    required this.lastMessage,
    required this.updatedAt,
    this.messages = const [],
  });

  factory Conversation.fromJson(Map<String, dynamic> json) {
    return Conversation(
      id: json['id'] ?? '',
      title: json['title'] ?? '新对话',
      lastMessage: json['last_message'] ?? '',
      updatedAt: DateTime.tryParse(json['updated_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class ChatProvider extends ChangeNotifier {
  final ApiService _api;
  final WebSocketService _ws;
  
  List<Conversation> _conversations = [];
  Conversation? _currentConversation;
  List<Message> _messages = [];
  bool _isLoading = false;
  String? _error;

  ChatProvider() : _api = ApiService(), _ws = WebSocketService() {
    _ws.messages.listen(_onWebSocketMessage);
  }

  List<Conversation> get conversations => _conversations;
  Conversation? get currentConversation => _currentConversation;
  List<Message> get messages => _messages;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void connectWebSocket(String serverUrl, String userId) {
    _ws.connect(serverUrl, userId);
  }

  void _onWebSocketMessage(Map<String, dynamic> data) {
    final message = Message.fromJson(data);
    if (_currentConversation?.id == data['conversation_id']) {
      _messages.add(message);
      notifyListeners();
    }
  }

  Future<void> loadConversations(String userId) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final data = await _api.getConversations(userId);
      _conversations = data.map((e) => Conversation.fromJson(e)).toList();
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> createConversation(String userId, String? agentId) async {
    _isLoading = true;
    notifyListeners();

    try {
      final data = await _api.createConversation({
        'user_id': userId,
        'agent_id': agentId ?? 'default',
      });
      if (data != null) {
        final conv = Conversation.fromJson(data);
        _conversations.insert(0, conv);
        _currentConversation = conv;
        _messages = [];
      }
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> selectConversation(Conversation conv) async {
    _currentConversation = conv;
    await loadMessages(conv.id);
  }

  Future<void> loadMessages(String conversationId) async {
    _isLoading = true;
    notifyListeners();

    try {
      final data = await _api.getMessages(conversationId);
      _messages = data.map((e) => Message.fromJson(e)).toList();
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> sendMessage(String conversationId, String content) async {
    // 先添加用户消息到UI
    final userMessage = Message(
      id: 'temp-${DateTime.now().millisecondsSinceEpoch}',
      type: 'text',
      content: content,
      sender: 'user',
      senderId: 'user',
      createdAt: DateTime.now(),
    );
    _messages.add(userMessage);
    notifyListeners();

    try {
      final data = await _api.sendMessage({
        'conversation_id': conversationId,
        'type': 'text',
        'content': content,
        'sender': 'user',
        'sender_id': 'user',
      });
      
      if (data != null) {
        final botMessage = Message.fromJson(data);
        _messages.add(botMessage);
      }
    } catch (e) {
      _error = e.toString();
    }

    notifyListeners();
  }

  Future<void> deleteConversation(String conversationId) async {
    try {
      await _api.deleteConversation(conversationId);
      _conversations.removeWhere((c) => c.id == conversationId);
      if (_currentConversation?.id == conversationId) {
        _currentConversation = null;
        _messages = [];
      }
    } catch (e) {
      _error = e.toString();
    }
    notifyListeners();
  }

  @override
  void dispose() {
    _ws.dispose();
    super.dispose();
  }
}
