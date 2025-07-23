import 'dart:convert';
import 'package:http/http.dart' as http;

void main() async {
  print('=== Testing Gate.io Public API ===');

  // Test public endpoint first
  await testPublicAPI();

  print('\n=== Checking API Key Status ===');

  // Note: These API credentials might be invalid or expired
  // The pattern suggests they could be demo/test keys
  const currentKey = '472667d2eebf07b206f8599bf35bbcb1';
  const currentSecret =
      'e24966020b7241e521966b87deb401d6dac0f90a1ccbf3f85db54ecc09faf33f';

  print('Current API Key: $currentKey');
  print('API Key length: ${currentKey.length}');
  print('Secret length: ${currentSecret.length}');

  // Check if the key format looks valid
  if (currentKey.length == 32 && currentSecret.length == 64) {
    print('✅ Key and secret have correct lengths');
  } else {
    print('❌ Key or secret length seems incorrect');
    print('Expected: Key=32 chars, Secret=64 chars');
    print(
      'Actual: Key=${currentKey.length} chars, Secret=${currentSecret.length} chars',
    );
  }

  // Check if they're hex strings
  final keyIsHex = RegExp(r'^[a-f0-9]+$').hasMatch(currentKey);
  final secretIsHex = RegExp(r'^[a-f0-9]+$').hasMatch(currentSecret);

  print('Key is hex: $keyIsHex');
  print('Secret is hex: $secretIsHex');

  if (!keyIsHex || !secretIsHex) {
    print('⚠️  Keys might not be in correct format (should be hex)');
  }
}

Future<void> testPublicAPI() async {
  try {
    print('🌐 Testing public API endpoint...');

    final response = await http.get(
      Uri.parse('https://api.gateio.ws/api/v4/spot/currencies'),
      headers: {'Accept': 'application/json'},
    );

    print('Status: ${response.statusCode}');

    if (response.statusCode == 200) {
      final data = json.decode(response.body) as List;
      print('✅ Public API works! Found ${data.length} currencies');

      // Show first few currencies
      for (int i = 0; i < 3 && i < data.length; i++) {
        final currency = data[i];
        print(
          'Currency ${i + 1}: ${currency['currency']} - ${currency['delisted'] == false ? 'Active' : 'Delisted'}',
        );
      }
    } else {
      print('❌ Public API failed: ${response.statusCode} - ${response.body}');
    }
  } catch (e) {
    print('❌ Public API error: $e');
  }
}
