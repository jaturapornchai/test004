import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:http/http.dart' as http;

void main() async {
  print('=== Advanced Gate.io API Debug ===');

  const apiKey = '472667d2eebf07b206f8599bf35bbcb1';
  const apiSecret =
      'e24966020b7241e521966b87deb401d6dac0f90a1ccbf3f85db54ecc09faf33f';

  await testMultipleSignatureMethods(apiKey, apiSecret);
}

Future<void> testMultipleSignatureMethods(
  String apiKey,
  String apiSecret,
) async {
  final baseUrl = 'https://api.gateio.ws/api/v4';
  final method = 'GET';
  final path = '/spot/accounts';
  final query = '';
  final body = '';

  // Test multiple timestamp and signature methods
  final testCases = [
    {
      'name': 'Integer timestamp (Go style)',
      'timestamp': (DateTime.now().millisecondsSinceEpoch / 1000)
          .floor()
          .toString(),
    },
    {
      'name': 'Float timestamp (Python style)',
      'timestamp': (DateTime.now().millisecondsSinceEpoch / 1000.0).toString(),
    },
    {
      'name': 'Millisecond timestamp',
      'timestamp': DateTime.now().millisecondsSinceEpoch.toString(),
    },
  ];

  for (var testCase in testCases) {
    print('\n🧪 Testing: ${testCase['name']}');
    print('Timestamp: ${testCase['timestamp']}');

    // Hash body
    final bodyHash = sha512.convert(utf8.encode(body)).toString();

    // Create message string
    final message =
        '$method\n$path\n$query\n$bodyHash\n${testCase['timestamp']}';

    // Generate signature
    final hmac = Hmac(sha512, utf8.encode(apiSecret));
    final signature = hmac.convert(utf8.encode(message)).toString();

    print('Message: ${message.replaceAll('\n', '\\n')}');
    print('Signature: $signature');

    // Make request
    final headers = {
      'KEY': apiKey,
      'SIGN': signature,
      'Timestamp': testCase['timestamp']!,
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };

    try {
      final response = await http.get(
        Uri.parse('$baseUrl$path'),
        headers: headers,
      );

      print('Status: ${response.statusCode}');
      if (response.statusCode == 200) {
        print('✅ SUCCESS!');
        final data = json.decode(response.body) as List;
        print('Found ${data.length} balances');
        break; // Success, stop testing
      } else {
        print('❌ Failed: ${response.body}');
      }
    } catch (e) {
      print('❌ Error: $e');
    }
  }

  // Test alternative empty string handling
  print('\n🧪 Testing alternative empty query handling...');
  final timestamp = (DateTime.now().millisecondsSinceEpoch / 1000)
      .floor()
      .toString();
  final bodyHash = sha512.convert(utf8.encode('')).toString();

  // Try with different empty query representations
  final queryVariants = ['', ' ', null];

  for (var i = 0; i < queryVariants.length; i++) {
    final queryVar = queryVariants[i] ?? '';
    print('\nQuery variant $i: "${queryVar}"');

    final message = '$method\n$path\n$queryVar\n$bodyHash\n$timestamp';
    final hmac = Hmac(sha512, utf8.encode(apiSecret));
    final signature = hmac.convert(utf8.encode(message)).toString();

    print('Message: ${message.replaceAll('\n', '\\n')}');

    final headers = {
      'KEY': apiKey,
      'SIGN': signature,
      'Timestamp': timestamp,
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };

    try {
      final response = await http.get(
        Uri.parse('$baseUrl$path'),
        headers: headers,
      );

      print('Status: ${response.statusCode}');
      if (response.statusCode == 200) {
        print('✅ SUCCESS with query variant $i!');
        break;
      } else {
        print('❌ Failed: ${response.body}');
      }
    } catch (e) {
      print('❌ Error: $e');
    }
  }

  // Check if credentials might be wrong by trying a simple request
  print('\n🧪 Testing if credentials are valid...');
  print('API Key: $apiKey');
  print('API Secret length: ${apiSecret.length} characters');
  print('API Secret (first 10 chars): ${apiSecret.substring(0, 10)}...');
}
