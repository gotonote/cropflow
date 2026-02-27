import 'package:dio/dio.dart';

class ApiService {
  static const String baseUrl = 'http://localhost:8080/api';
  
  late final Dio _dio;
  
  ApiService() {
    _dio = Dio(BaseOptions(
      baseUrl: baseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 30),
      headers: {
        'Content-Type': 'application/json',
      },
    ));
    
    _dio.interceptors.add(LogInterceptor(
      requestBody: true,
      responseBody: true,
    ));
  }

  // ========== 智能体 ==========
  
  Future<List<dynamic>> getAgents() async {
    try {
      final response = await _dio.get('/agents');
      return response.data;
    } catch (e) {
      return [];
    }
  }

  Future<Map<String, dynamic>?> getAgent(int id) async {
    try {
      final response = await _dio.get('/agents/$id');
      return response.data;
    } catch (e) {
      return null;
    }
  }

  Future<Map<String, dynamic>?> createAgent(Map<String, dynamic> data) async {
    try {
      final response = await _dio.post('/agents', data: data);
      return response.data;
    } catch (e) {
      return null;
    }
  }

  Future<bool> updateAgent(int id, Map<String, dynamic> data) async {
    try {
      await _dio.put('/agents/$id', data: data);
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<bool> deleteAgent(int id) async {
    try {
      await _dio.delete('/agents/$id');
      return true;
    } catch (e) {
      return false;
    }
  }

  // ========== 流程 ==========
  
  Future<List<dynamic>> getFlows() async {
    try {
      final response = await _dio.get('/flows');
      return response.data;
    } catch (e) {
      return [];
    }
  }

  Future<Map<String, dynamic>?> createFlow(Map<String, dynamic> data) async {
    try {
      final response = await _dio.post('/flows', data: data);
      return response.data;
    } catch (e) {
      return null;
    }
  }

  Future<bool> updateFlow(int id, Map<String, dynamic> data) async {
    try {
      await _dio.put('/flows/$id', data: data);
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<bool> deleteFlow(int id) async {
    try {
      await _dio.delete('/flows/$id');
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<String?> executeFlow(int id, Map<String, dynamic> data) async {
    try {
      final response = await _dio.post('/flows/$id/execute', data: data);
      return response.data['output'];
    } catch (e) {
      return null;
    }
  }

  // ========== 渠道 ==========
  
  Future<List<dynamic>> getChannels() async {
    try {
      final response = await _dio.get('/channels');
      return response.data;
    } catch (e) {
      return [];
    }
  }

  Future<Map<String, dynamic>?> createChannel(Map<String, dynamic> data) async {
    try {
      final response = await _dio.post('/channels', data: data);
      return response.data;
    } catch (e) {
      return null;
    }
  }

  Future<bool> updateChannel(int id, Map<String, dynamic> data) async {
    try {
      await _dio.put('/channels/$id', data: data);
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<bool> deleteChannel(int id) async {
    try {
      await _dio.delete('/channels/$id');
      return true;
    } catch (e) {
      return false;
    }
  }

  // ========== 聊天 ==========
  
  Future<List<dynamic>> getConversations(String userId) async {
    try {
      final response = await _dio.get('/chat/conversations', queryParameters: {'user_id': userId});
      return response.data;
    } catch (e) {
      return [];
    }
  }

  Future<Map<String, dynamic>?> createConversation(Map<String, dynamic> data) async {
    try {
      final response = await _dio.post('/chat/conversations', data: data);
      return response.data;
    } catch (e) {
      return null;
    }
  }

  Future<List<dynamic>> getMessages(String conversationId) async {
    try {
      final response = await _dio.get('/chat/conversations/$conversationId/messages');
      return response.data;
    } catch (e) {
      return [];
    }
  }

  Future<Map<String, dynamic>?> sendMessage(Map<String, dynamic> data) async {
    try {
      final response = await _dio.post('/chat/messages', data: data);
      return response.data;
    } catch (e) {
      return null;
    }
  }

  // ========== 配置 ==========
  
  void setBaseUrl(String url) {
    _dio.options.baseUrl = url;
  }

  void setToken(String token) {
    _dio.options.headers['Authorization'] = 'Bearer $token';
  }
}
