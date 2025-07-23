import 'dart:convert';
import 'dart:io';
import 'package:crypto/crypto.dart';
import 'package:http/http.dart' as http;

class GateApiClientNew {
  static const String baseUrl = 'https://api.gateio.ws/api/v4';
  final String apiKey;
  final String apiSecret;

  GateApiClientNew({required this.apiKey, required this.apiSecret});

  String _generateSignature(
    String method,
    String url,
    String queryString,
    String body,
  ) {
    // Get current timestamp in seconds (integer like Go client)
    final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;

    // Hash the body with SHA512
    final bodyBytes = utf8.encode(body);
    final hashedPayload = sha512.convert(bodyBytes).toString();

    // Create the signature string exactly as shown in official Go client
    final signString =
        '$method\n$url\n$queryString\n$hashedPayload\n$timestamp';

    print('🔐 Gate.io Signature Debug:');
    print('Timestamp: $timestamp');
    print('Method: $method');
    print('URL: $url');
    print('Query: $queryString');
    print('Body: $body');
    print('Body Hash: $hashedPayload');
    print('Sign String: ${signString.replaceAll('\n', '\\n')}');

    // Generate HMAC-SHA512 signature
    final key = utf8.encode(apiSecret);
    final messageBytes = utf8.encode(signString);
    final hmac = Hmac(sha512, key);
    final signature = hmac.convert(messageBytes).toString();

    print('Signature: $signature');

    return signature;
  }

  Map<String, String> _createHeaders(
    String method,
    String url,
    String queryString,
    String body,
  ) {
    final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
    final signature = _generateSignature(method, url, queryString, body);

    return {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'KEY': apiKey,
      'Timestamp': timestamp.toString(),
      'SIGN': signature,
    };
  }

  Future<List<Map<String, dynamic>>> getSpotAccounts() async {
    const method = 'GET';
    const path = '/spot/accounts';
    const queryString = '';
    const body = '';

    try {
      final headers = _createHeaders(method, path, queryString, body);
      final url = '$baseUrl$path';

      print('🌐 Gate.io Spot Request:');
      print('URL: $url');
      print('Headers: $headers');

      final response = await http
          .get(Uri.parse(url), headers: headers)
          .timeout(const Duration(seconds: 30));

      print('📥 Gate.io Spot Response:');
      print('Status: ${response.statusCode}');
      print('Body length: ${response.body.length}');

      if (response.statusCode == 200) {
        final List<dynamic> data = jsonDecode(response.body);
        return data.cast<Map<String, dynamic>>();
      } else {
        print(
          '❌ Gate.io Spot Error: ${response.statusCode} - ${response.body}',
        );
        throw Exception(
          'Failed to load spot accounts: ${response.statusCode} - ${response.body}',
        );
      }
    } catch (e) {
      print('❌ Gate.io Spot Network Error: $e');
      rethrow;
    }
  }

  Future<Map<String, dynamic>?> getFuturesAccount() async {
    const method = 'GET';
    const path = '/futures/usdt/accounts';
    const queryString = '';
    const body = '';

    try {
      final headers = _createHeaders(method, path, queryString, body);
      final url = '$baseUrl$path';

      print('🌐 Gate.io Futures Request:');
      print('URL: $url');
      print('Headers: $headers');

      final response = await http
          .get(Uri.parse(url), headers: headers)
          .timeout(const Duration(seconds: 30));

      print('📥 Gate.io Futures Response:');
      print('Status: ${response.statusCode}');
      print('Body: ${response.body}');

      if (response.statusCode == 200) {
        return jsonDecode(response.body) as Map<String, dynamic>;
      } else {
        print(
          '❌ Gate.io Futures Error: ${response.statusCode} - ${response.body}',
        );
        return null;
      }
    } catch (e) {
      print('❌ Gate.io Futures Network Error: $e');
      return null;
    }
  }
}
