import 'dart:convert';
import 'dart:io';
import 'package:crypto/crypto.dart';

class GateApiClient {
  static const String baseUrl = 'https://api.gateio.ws/api/v4';
  final String apiKey;
  final String apiSecret;

  GateApiClient({required this.apiKey, required this.apiSecret});

  String _createSignature(
    String method,
    String url,
    String query,
    String body,
    String timestamp,
  ) {
    // Hash the request body with SHA512 first (this was the missing piece!)
    final hashedBody = sha512.convert(utf8.encode(body)).toString();
    // Use full path including /api/v4 for signature
    final fullPath = '/api/v4$url';
    final message = '$method\n$fullPath\n$query\n$hashedBody\n$timestamp';
    var key = utf8.encode(apiSecret);
    var bytes = utf8.encode(message);
    var hmacSha512 = Hmac(sha512, key);
    var digest = hmacSha512.convert(bytes);
    return digest.toString();
  }

  // Public method to generate signature (for external use)
  String generateSignature(
    String method,
    String url,
    String query,
    String body,
  ) {
    final timestamp = (DateTime.now().millisecondsSinceEpoch / 1000)
        .round()
        .toString();
    return _createSignature(method, url, query, body, timestamp);
  }

  Future<List<dynamic>> getSpotAccounts() async {
    try {
      const method = 'GET';
      const url = '/spot/accounts';
      const query = '';
      const body = '';
      final timestamp = (DateTime.now().millisecondsSinceEpoch / 1000)
          .round()
          .toString();

      final signature = _createSignature(method, url, query, body, timestamp);

      final request = await HttpClient().getUrl(Uri.parse('$baseUrl$url'));
      request.headers.set('Accept', 'application/json');
      request.headers.set('Content-Type', 'application/json');
      request.headers.set('KEY', apiKey);
      request.headers.set('Timestamp', timestamp);
      request.headers.set('SIGN', signature);

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      if (response.statusCode == 200) {
        return jsonDecode(responseBody) as List;
      } else {
        print('Gate.io Spot API Error: ${response.statusCode} - $responseBody');
        throw Exception(
          'Failed to load spot accounts: ${response.statusCode} - $responseBody',
        );
      }
    } catch (e) {
      print('Gate.io Spot Network Error: $e');
      throw Exception('Network error: $e');
    }
  }

  Future<Map<String, dynamic>?> getFuturesAccount() async {
    try {
      const method = 'GET';
      const url = '/futures/usdt/accounts';
      const query = '';
      const body = '';
      final timestamp = (DateTime.now().millisecondsSinceEpoch / 1000)
          .round()
          .toString();

      final signature = _createSignature(method, url, query, body, timestamp);

      final request = await HttpClient().getUrl(Uri.parse('$baseUrl$url'));
      request.headers.set('Accept', 'application/json');
      request.headers.set('Content-Type', 'application/json');
      request.headers.set('KEY', apiKey);
      request.headers.set('Timestamp', timestamp);
      request.headers.set('SIGN', signature);

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      if (response.statusCode == 200) {
        return jsonDecode(responseBody);
      } else {
        print(
          'Gate.io Futures API Error: ${response.statusCode} - $responseBody',
        );
        return null;
      }
    } catch (e) {
      print('Gate.io Futures Network Error: $e');
      return null;
    }
  }

  Future<List<Map<String, dynamic>>> getAllTickers() async {
    try {
      final request = await HttpClient().getUrl(
        Uri.parse('$baseUrl/spot/tickers'),
      );
      request.headers.set('Accept', 'application/json');

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      if (response.statusCode == 200) {
        return List<Map<String, dynamic>>.from(jsonDecode(responseBody));
      } else {
        throw Exception('Failed to load tickers: ${response.statusCode}');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  // Get today's spot trade history
  Future<List<Map<String, dynamic>>> getTodaySpotTradeHistory() async {
    try {
      const method = 'GET';
      const url = '/spot/my_trades';

      final now = DateTime.now();
      final todayStart = DateTime(now.year, now.month, now.day);
      final startTime = (todayStart.millisecondsSinceEpoch / 1000).round();
      final endTime = (now.millisecondsSinceEpoch / 1000).round();

      final query = 'from=$startTime&to=$endTime';
      const body = '';
      final timestamp = (DateTime.now().millisecondsSinceEpoch / 1000)
          .round()
          .toString();
      final signature = _createSignature(method, url, query, body, timestamp);

      final request = await HttpClient().getUrl(
        Uri.parse('$baseUrl$url?$query'),
      );
      request.headers.set('Accept', 'application/json');
      request.headers.set('KEY', apiKey);
      request.headers.set('Timestamp', timestamp);
      request.headers.set('SIGN', signature);

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      if (response.statusCode == 200) {
        return List<Map<String, dynamic>>.from(jsonDecode(responseBody));
      } else {
        throw Exception(
          'Failed to load spot trade history: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  // Get today's futures trade history
  Future<List<Map<String, dynamic>>> getTodayFuturesTradeHistory() async {
    try {
      const method = 'GET';
      const url = '/futures/usdt/my_trades';

      final now = DateTime.now();
      final todayStart = DateTime(now.year, now.month, now.day);
      final startTime = (todayStart.millisecondsSinceEpoch / 1000).round();
      final endTime = (now.millisecondsSinceEpoch / 1000).round();

      final query = 'from=$startTime&to=$endTime';
      const body = '';
      final timestamp = (DateTime.now().millisecondsSinceEpoch / 1000)
          .round()
          .toString();
      final signature = _createSignature(method, url, query, body, timestamp);

      final request = await HttpClient().getUrl(
        Uri.parse('$baseUrl$url?$query'),
      );
      request.headers.set('Accept', 'application/json');
      request.headers.set('KEY', apiKey);
      request.headers.set('Timestamp', timestamp);
      request.headers.set('SIGN', signature);

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      if (response.statusCode == 200) {
        return List<Map<String, dynamic>>.from(jsonDecode(responseBody));
      } else {
        throw Exception(
          'Failed to load futures trade history: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }
}
