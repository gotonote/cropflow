import 'dart:async';
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';

class WebSocketService {
  WebSocketChannel? _channel;
  final _messageController = StreamController<Map<String, dynamic>>.broadcast();
  final _connectionController = StreamController<bool>.broadcast();
  
  String? _serverUrl;
  String? _userId;
  bool _isConnected = false;
  Timer? _reconnectTimer;
  Timer? _pingTimer;

  Stream<Map<String, dynamic>> get messages => _messageController.stream;
  Stream<bool> get connectionStatus => _connectionController.stream;
  bool get isConnected => _isConnected;

  void connect(String serverUrl, String userId) {
    _serverUrl = serverUrl;
    _userId = userId;
    _doConnect();
  }

  void _doConnect() {
    if (_serverUrl == null || _userId == null) return;

    try {
      final uri = Uri.parse('${_serverUrl.replaceAll('http', 'ws')}/ws?user_id=$_userId');
      _channel = WebSocketChannel.connect(uri);

      _channel!.stream.listen(
        _onMessage,
        onError: _onError,
        onDone: _onDone,
      );

      _isConnected = true;
      _connectionController.add(true);
      _startPing();
    } catch (e) {
      _isConnected = false;
      _connectionController.add(false);
      _scheduleReconnect();
    }
  }

  void _onMessage(dynamic data) {
    try {
      final message = jsonDecode(data as String) as Map<String, dynamic>;
      _messageController.add(message);
    } catch (e) {
      // 忽略解析错误
    }
  }

  void _onError(dynamic error) {
    _isConnected = false;
    _connectionController.add(false);
    _scheduleReconnect();
  }

  void _onDone() {
    _isConnected = false;
    _connectionController.add(false);
    _scheduleReconnect();
  }

  void _scheduleReconnect() {
    _pingTimer?.cancel();
    _reconnectTimer?.cancel();
    _reconnectTimer = Timer(const Duration(seconds: 3), () {
      _doConnect();
    });
  }

  void _startPing() {
    _pingTimer?.cancel();
    _pingTimer = Timer.periodic(const Duration(seconds: 30), (_) {
      send({'type': 'ping'});
    });
  }

  void send(Map<String, dynamic> message) {
    if (_channel != null && _isConnected) {
      _channel!.sink.add(jsonEncode(message));
    }
  }

  void sendMessage(String conversationId, String content) {
    send({
      'type': 'message',
      'conversation_id': conversationId,
      'content': content,
    });
  }

  void disconnect() {
    _pingTimer?.cancel();
    _reconnectTimer?.cancel();
    _channel?.sink.close();
    _channel = null;
    _isConnected = false;
    _connectionController.add(false);
  }

  void dispose() {
    disconnect();
    _messageController.close();
    _connectionController.close();
  }
}
