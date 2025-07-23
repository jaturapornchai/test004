import 'dart:convert';
import 'dart:io';
import 'package:crypto/crypto.dart';
import 'package:http/http.dart' as http;

void main() async {
  print('=== Testing Gate.io API with Go-style Authentication ===');

  const apiKey = '472667d2eebf07b206f8599bf35bbcb1';
  const apiSecret =
      'e24966020b7241e521966b87deb401d6dac0f90a1ccbf3f85db54ecc09faf33f';
  const baseUrl = 'https://api.gateio.ws/api/v4';

  await testGateIOAuth(apiKey, apiSecret, baseUrl);
}

Future<void> testGateIOAuth(
  String apiKey,
  String apiSecret,
  String baseUrl,
) async {
  try {
    print('🟢 Testing Gate.io Spot Account...');

    // Create request components exactly like Go client
    final method = 'GET';
    final path = '/spot/accounts';
    final queryString = '';
    final body = '';

    // Generate timestamp as integer (like Go client)
    final timestamp = (DateTime.now().millisecondsSinceEpoch / 1000).floor();

    // Hash the body with SHA512 (empty body for GET request)
    final bodyBytes = utf8.encode(body);
    final hashedPayload = sha512.convert(bodyBytes).toString();

    // Create the message string exactly like Go client:
    // method + '\n' + path + '\n' + query + '\n' + hashedPayload + '\n' + timestamp
    final message = '$method\n$path\n$queryString\n$hashedPayload\n$timestamp';

    // Generate HMAC-SHA512 signature
    final secretBytes = utf8.encode(apiSecret);
    final messageBytes = utf8.encode(message);
    final hmacSha512 = Hmac(sha512, secretBytes);
    final signature = hmacSha512.convert(messageBytes).toString();

    print('🔐 Gate.io Authentication Debug:');
    print('Method: $method');
    print('Path: $path');
    print('Query: "$queryString"');
    print('Body: "$body"');
    print('Hashed Payload: $hashedPayload');
    print('Timestamp: $timestamp');
    print('Message: ${message.replaceAll('\n', '\\n')}');
    print('Signature: $signature');
    print('');

    // Create headers exactly like Go client
    final headers = {
      'KEY': apiKey,
      'SIGN': signature,
      'Timestamp': timestamp.toString(),
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };

    print('🌐 Request Headers:');
    headers.forEach((key, value) {
      print('$key: $value');
    });
    print('');

    // Make the request
    final url = Uri.parse('$baseUrl$path');
    final response = await http.get(url, headers: headers);

    print('📥 Response:');
    print('Status: ${response.statusCode}');
    print('Body: ${response.body}');

    if (response.statusCode == 200) {
      final data = json.decode(response.body) as List;
      print('✅ Gate.io Spot Account Success!');
      print('Found ${data.length} balances');

      double total = 0.0;
      for (var account in data) {
        final currency = account['currency'] ?? '';
        final available = double.tryParse(account['available'] ?? '0') ?? 0.0;
        final locked = double.tryParse(account['locked'] ?? '0') ?? 0.0;
        final balance = available + locked;

        if (balance > 0) {
          print('$currency: $balance (available: $available, locked: $locked)');
          if (currency == 'USDT') {
            total += balance;
          }
        }
      }
      print('Total USDT balance: \$${total.toStringAsFixed(2)}');
    } else {
      print('❌ Gate.io Error: ${response.statusCode} - ${response.body}');
    }
  } catch (e) {
    print('❌ Gate.io Exception: $e');
  }
}
