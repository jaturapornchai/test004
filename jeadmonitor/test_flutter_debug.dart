import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:http/http.dart' as http;

Future<void> main() async {
  // ใส่ API credentials ตรงนี้
  const apiKey = '472667d2eebf07b206f8599bf35bbcb1';
  const apiSecret =
      'e24966020b7241e521966b87deb401d6dac0f90a1ccbf3f85db54ecc09faf33f';

  print('=== Flutter Debug Test ===');

  try {
    final client = GateApiDebugClient(apiKey: apiKey, apiSecret: apiSecret);

    // Test spot accounts
    print('\n1. Testing Spot Accounts...');
    final spotAccounts = await client.getSpotAccounts();
    print('Spot accounts count: ${spotAccounts.length}');

    // Test futures account
    print('\n2. Testing Futures Account...');
    final futuresAccount = await client.getFuturesAccount();
    if (futuresAccount != null) {
      print('Futures Total: \$${futuresAccount['total']}');
    }
  } catch (e) {
    print('Error: $e');
  }
}

class GateApiDebugClient {
  static const String baseUrl = 'https://api.gateio.ws/api/v4';
  final String apiKey;
  final String apiSecret;

  GateApiDebugClient({required this.apiKey, required this.apiSecret});

  String _createSignature(
    String method,
    String url,
    String query,
    String body,
    String timestamp,
  ) {
    // Hash the request body with SHA512 first
    final hashedBody = sha512.convert(utf8.encode(body)).toString();
    // Use full path including /api/v4 for signature
    final fullPath = '/api/v4$url';
    final message = '$method\n$fullPath\n$query\n$hashedBody\n$timestamp';

    print('\n=== FLUTTER REQUEST DEBUG ===');
    print('Method: $method');
    print('URL: $url');
    print('Full Path: $fullPath');
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
    const method = 'GET';
    const url = '/spot/accounts';
    const query = '';
    const body = '';
    final timestamp = (DateTime.now().millisecondsSinceEpoch / 1000)
        .round()
        .toString();

    final signature = _createSignature(method, url, query, body, timestamp);

    final response = await http.get(
      Uri.parse('$baseUrl$url'),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'KEY': apiKey,
        'Timestamp': timestamp,
        'SIGN': signature,
      },
    );

    print('\nFLUTTER RESPONSE:');
    print('Status: ${response.statusCode}');
    print('Headers: ${response.headers}');
    print('Body: ${response.body}');

    if (response.statusCode == 200) {
      return jsonDecode(response.body) as List;
    } else {
      throw Exception('Failed: ${response.statusCode} - ${response.body}');
    }
  }

  Future<Map<String, dynamic>?> getFuturesAccount() async {
    const method = 'GET';
    const url = '/futures/usdt/accounts';
    const query = '';
    const body = '';
    final timestamp = (DateTime.now().millisecondsSinceEpoch / 1000)
        .round()
        .toString();

    final signature = _createSignature(method, url, query, body, timestamp);

    final response = await http.get(
      Uri.parse('$baseUrl$url'),
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'KEY': apiKey,
        'Timestamp': timestamp,
        'SIGN': signature,
      },
    );

    print('\nFUTURES RESPONSE:');
    print('Status: ${response.statusCode}');
    print('Body: ${response.body}');

    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    } else {
      print('Futures failed: ${response.statusCode} - ${response.body}');
      return null;
    }
  }
}
