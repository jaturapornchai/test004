import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:http/http.dart' as http;

class GateApiClientDebug {
  static const String baseUrl = 'https://api.gateio.ws/api/v4';
  final String apiKey;
  final String apiSecret;

  GateApiClientDebug({required this.apiKey, required this.apiSecret});

  String _createSignature(
    String method,
    String url,
    String query,
    String body,
    String timestamp,
  ) {
    // Hash the request body with SHA512 first
    final hashedBody = sha512.convert(utf8.encode(body)).toString();
    final message = '$method\n$url\n$query\n$hashedBody\n$timestamp';

    print('=== Gate.io Flutter Debug ===');
    print('Method: $method');
    print('URL: $url');
    print('Query: $query');
    print('Body: $body');
    print('Body Hash: $hashedBody');
    print('Timestamp: $timestamp');
    print('Message: "$message"');

    var key = utf8.encode(apiSecret);
    var bytes = utf8.encode(message);
    var hmacSha512 = Hmac(sha512, key);
    var digest = hmacSha512.convert(bytes);

    print('Signature: $digest');
    print('API Key: $apiKey');
    print('=============================');

    return digest.toString();
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

      final response = await http
          .get(
            Uri.parse('$baseUrl$url'),
            headers: {
              'Accept': 'application/json',
              'Content-Type': 'application/json',
              'KEY': apiKey,
              'Timestamp': timestamp,
              'SIGN': signature,
            },
          )
          .timeout(const Duration(seconds: 15));

      print('Response Status: ${response.statusCode}');
      print('Response Headers: ${response.headers}');
      print('Response Body: ${response.body}');

      if (response.statusCode == 200) {
        return jsonDecode(response.body) as List;
      } else {
        print(
          'Gate.io Spot API Error: ${response.statusCode} - ${response.body}',
        );
        throw Exception(
          'Failed to load spot accounts: ${response.statusCode} - ${response.body}',
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

      final response = await http
          .get(
            Uri.parse('$baseUrl$url'),
            headers: {
              'Accept': 'application/json',
              'Content-Type': 'application/json',
              'KEY': apiKey,
              'Timestamp': timestamp,
              'SIGN': signature,
            },
          )
          .timeout(const Duration(seconds: 15));

      print('Futures Response Status: ${response.statusCode}');
      print('Futures Response Body: ${response.body}');

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        print(
          'Gate.io Futures API Error: ${response.statusCode} - ${response.body}',
        );
        return null;
      }
    } catch (e) {
      print('Gate.io Futures Network Error: $e');
      return null;
    }
  }
}
