import 'package:dio/dio.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class ApiClient {
  late final Dio _dio;
  final String baseUrl;
  String? _token;
  WebSocketChannel? _wsChannel;

  ApiClient({required this.baseUrl}) {
    _dio = Dio(BaseOptions(
      baseUrl: baseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 30),
      headers: {
        'Content-Type': 'application/json',
      },
    ));

    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) {
        if (_token != null) {
          options.headers['Authorization'] = 'Bearer $_token';
        }
        return handler.next(options);
      },
      onError: (error, handler) {
        // 统一错误处理
        return handler.next(error);
      },
    ));
  }

  void setToken(String token) {
    _token = token;
  }

  void clearToken() {
    _token = null;
  }

  // ========== 认证 ==========
  
  Future<Map<String, dynamic>> login(String username, String password) async {
    final response = await _dio.post('/api/auth/login', data: {
      'username': username,
      'password': password,
    });
    return response.data;
  }

  Future<Map<String, dynamic>> register(String username, String password, String email) async {
    final response = await _dio.post('/api/auth/register', data: {
      'username': username,
      'password': password,
      'email': email,
    });
    return response.data;
  }

  // ========== 智能体 ==========

  Future<List<dynamic>> getAgents() async {
    final response = await _dio.get('/api/agents');
    return response.data;
  }

  Future<Map<String, dynamic>> getAgent(String id) async {
    final response = await _dio.get('/api/agents/$id');
    return response.data;
  }

  Future<Map<String, dynamic>> createAgent(Map<String, dynamic> data) async {
    final response = await _dio.post('/api/agents', data: data);
    return response.data;
  }

  Future<Map<String, dynamic>> updateAgent(String id, Map<String, dynamic> data) async {
    final response = await _dio.put('/api/agents/$id', data: data);
    return response.data;
  }

  Future<void> deleteAgent(String id) async {
    await _dio.delete('/api/agents/$id');
  }

  // ========== 流程 ==========

  Future<List<dynamic>> getFlows() async {
    final response = await _dio.get('/api/flows');
    return response.data;
  }

  Future<Map<String, dynamic>> getFlow(String id) async {
    final response = await _dio.get('/api/flows/$id');
    return response.data;
  }

  Future<Map<String, dynamic>> createFlow(Map<String, dynamic> data) async {
    final response = await _dio.post('/api/flows', data: data);
    return response.data;
  }

  Future<Map<String, dynamic>> updateFlow(String id, Map<String, dynamic> data) async {
    final response = await _dio.put('/api/flows/$id', data: data);
    return response.data;
  }

  Future<void> deleteFlow(String id) async {
    await _dio.delete('/api/flows/$id');
  }

  Future<Map<String, dynamic>> executeFlow(String id, Map<String, dynamic> input) async {
    final response = await _dio.post('/api/flows/$id/execute', data: input);
    return response.data;
  }

  // ========== 会话 ==========

  Future<List<dynamic>> getConversations() async {
    final response = await _dio.get('/api/chat/conversations');
    return response.data;
  }

  Future<Map<String, dynamic>> createConversation(Map<String, dynamic> data) async {
    final response = await _dio.post('/api/chat/conversations', data: data);
    return response.data;
  }

  Future<Map<String, dynamic>> getConversation(String id) async {
    final response = await _dio.get('/api/chat/conversations/$id');
    return response.data;
  }

  Future<void> deleteConversation(String id) async {
    await _dio.delete('/api/chat/conversations/$id');
  }

  // ========== 消息 ==========

  Future<List<dynamic>> getMessages(String convId, {int limit = 50, int offset = 0}) async {
    final response = await _dio.get('/api/chat/conversations/$convId/messages', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    return response.data;
  }

  Future<Map<String, dynamic>> sendMessage(String convId, Map<String, dynamic> data) async {
    final response = await _dio.post('/api/chat/conversations/$convId/messages', data: data);
    return response.data;
  }

  // ========== 渠道 ==========

  Future<List<dynamic>> getChannels() async {
    final response = await _dio.get('/api/channels');
    return response.data;
  }

  Future<Map<String, dynamic>> createChannel(Map<String, dynamic> data) async {
    final response = await _dio.post('/api/channels', data: data);
    return response.data;
  }

  Future<void> deleteChannel(String id) async {
    await _dio.delete('/api/channels/$id');
  }

  // ========== WebSocket ==========

  WebSocketChannel connectWebSocket(String userId, String conversationId) {
    final uri = Uri.parse('${baseUrl.replaceFirst('http', 'ws')}/ws?user_id=$userId&conversation_id=$conversationId');
    _wsChannel = WebSocketChannel.connect(uri);
    return _wsChannel!;
  }

  void disconnectWebSocket() {
    _wsChannel?.sink.close();
    _wsChannel = null;
  }
}
