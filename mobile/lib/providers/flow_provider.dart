import 'package:flutter/foundation.dart';
import '../services/api_service.dart';

class FlowNode {
  final String id;
  final String type;
  final Map<String, dynamic> data;
  final Map<String, dynamic> position;

  FlowNode({
    required this.id,
    required this.type,
    required this.data,
    required this.position,
  });

  factory FlowNode.fromJson(Map<String, dynamic> json) {
    return FlowNode(
      id: json['id'] ?? '',
      type: json['type'] ?? 'agent',
      data: json['data'] ?? {},
      position: json['position'] ?? {'x': 0, 'y': 0},
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'type': type,
    'data': data,
    'position': position,
  };
}

class FlowEdge {
  final String id;
  final String source;
  final String target;

  FlowEdge({
    required this.id,
    required this.source,
    required this.target,
  });

  factory FlowEdge.fromJson(Map<String, dynamic> json) {
    return FlowEdge(
      id: json['id'] ?? '',
      source: json['source'] ?? '',
      target: json['target'] ?? '',
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'source': source,
    'target': target,
  };
}

class Flow {
  final int id;
  final String name;
  final List<FlowNode> nodes;
  final List<FlowEdge> edges;
  final bool enabled;
  final DateTime updatedAt;

  Flow({
    required this.id,
    required this.name,
    required this.nodes,
    required this.edges,
    required this.enabled,
    required this.updatedAt,
  });

  factory Flow.fromJson(Map<String, dynamic> json) {
    return Flow(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      nodes: (json['nodes'] as List<dynamic>?)
          ?.map((e) => FlowNode.fromJson(e))
          .toList() ?? [],
      edges: (json['edges'] as List<dynamic>?)
          ?.map((e) => FlowEdge.fromJson(e))
          .toList() ?? [],
      enabled: json['enabled'] ?? true,
      updatedAt: DateTime.tryParse(json['updated_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class FlowProvider extends ChangeNotifier {
  final ApiService _api;
  
  List<Flow> _flows = [];
  Flow? _currentFlow;
  bool _isLoading = false;
  String? _error;
  String? _executionResult;

  FlowProvider() : _api = ApiService();

  List<Flow> get flows => _flows;
  Flow? get currentFlow => _currentFlow;
  bool get isLoading => _isLoading;
  String? get error => _error;
  String? get executionResult => _executionResult;

  Future<void> loadFlows() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final data = await _api.getFlows();
      _flows = data.map((e) => Flow.fromJson(e)).toList();
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> createFlow(String name) async {
    _isLoading = true;
    notifyListeners();

    try {
      final data = await _api.createFlow({
        'name': name,
        'nodes': [],
        'edges': [],
        'trigger_type': 'manual',
      });
      if (data != null) {
        _flows.insert(0, Flow.fromJson(data));
      }
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  void selectFlow(Flow flow) {
    _currentFlow = flow;
    _executionResult = null;
    notifyListeners();
  }

  Future<void> updateFlow(int id, List<FlowNode> nodes, List<FlowEdge> edges) async {
    try {
      await _api.updateFlow(id, {
        'nodes': nodes.map((e) => e.toJson()).toList(),
        'edges': edges.map((e) => e.toJson()).toList(),
      });
      
      final index = _flows.indexWhere((f) => f.id == id);
      if (index != -1) {
        _flows[index] = Flow(
          id: id,
          name: _flows[index].name,
          nodes: nodes,
          edges: edges,
          enabled: _flows[index].enabled,
          updatedAt: DateTime.now(),
        );
        if (_currentFlow?.id == id) {
          _currentFlow = _flows[index];
        }
      }
    } catch (e) {
      _error = e.toString();
    }
    notifyListeners();
  }

  Future<void> deleteFlow(int id) async {
    try {
      await _api.deleteFlow(id);
      _flows.removeWhere((f) => f.id == id);
      if (_currentFlow?.id == id) {
        _currentFlow = null;
      }
    } catch (e) {
      _error = e.toString();
    }
    notifyListeners();
  }

  Future<void> executeFlow(int id, String input) async {
    _isLoading = true;
    _executionResult = null;
    notifyListeners();

    try {
      _executionResult = await _api.executeFlow(id, {'input': input});
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }
}
