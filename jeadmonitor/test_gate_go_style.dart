import 'dart:convert';
import 'dart:io';
import 'package:crypto/crypto.dart';

class GateGoStyleTest {
  static const String baseUrl = 'https://api.gateio.ws';
  static const String apiKey = '472667d2eebf07b206f8599bf35bbcb1';
  static const String apiSecret =
      'e24966020b7241e521966b87deb401d6dac0f90a1ccbf3f85db54ecc09faf33f';

  // สร้าง signature แบบ Go style
  static String createSignature(
    String method,
    String path,
    String queryString,
    String body,
    int timestamp,
  ) {
    // สร้าง payload ตามรูปแบบ Go client
    String payload = '$method\n$path\n$queryString\n$body\n$timestamp';

    print('=== Go Style Signature Debug ===');
    print('Method: $method');
    print('Path: $path');
    print('Query String: $queryString');
    print('Body: $body');
    print('Timestamp: $timestamp');
    print('Payload: "$payload"');

    // สร้าง signature
    var secretBytes = utf8.encode(apiSecret);
    var payloadBytes = utf8.encode(payload);
    var hmac = Hmac(sha512, secretBytes);
    var signature = hmac.convert(payloadBytes);

    print('Signature: $signature');
    print('==============================');

    return signature.toString();
  }

  // Test Spot Account (Go style)
  static Future<void> testSpotAccount() async {
    print('\n🔥 Testing Spot Account (Go Style)...');

    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
      final method = 'GET';
      final path = '/api/v4/spot/accounts';
      final queryString = '';
      final body = '';

      final signature = createSignature(
        method,
        path,
        queryString,
        body,
        timestamp,
      );

      final client = HttpClient();
      final request = await client.getUrl(Uri.parse('$baseUrl$path'));

      // Headers ตามแบบ Go client
      request.headers.set('KEY', apiKey);
      request.headers.set('Timestamp', timestamp.toString());
      request.headers.set('SIGN', signature);
      request.headers.set('Content-Type', 'application/json');
      request.headers.set('Accept', 'application/json');

      print('Headers:');
      print('  KEY: $apiKey');
      print('  Timestamp: $timestamp');
      print('  SIGN: $signature');

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      print('Status Code: ${response.statusCode}');
      print('Response Body: $responseBody');

      if (response.statusCode == 200) {
        final data = jsonDecode(responseBody) as List;
        print('✅ Success! Found ${data.length} accounts');
        for (var account in data) {
          final currency = account['currency'];
          final available = account['available'];
          final locked = account['locked'];
          if (double.parse(available) > 0 || double.parse(locked) > 0) {
            print('  $currency: Available=$available, Locked=$locked');
          }
        }
      } else {
        print('❌ Failed: ${response.statusCode}');
      }

      client.close();
    } catch (e) {
      print('❌ Error: $e');
    }
  }

  // Test Futures Account (Go style)
  static Future<void> testFuturesAccount() async {
    print('\n🔥 Testing Futures Account (Go Style)...');

    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
      final method = 'GET';
      final path = '/api/v4/futures/usdt/accounts';
      final queryString = '';
      final body = '';

      final signature = createSignature(
        method,
        path,
        queryString,
        body,
        timestamp,
      );

      final client = HttpClient();
      final request = await client.getUrl(Uri.parse('$baseUrl$path'));

      request.headers.set('KEY', apiKey);
      request.headers.set('Timestamp', timestamp.toString());
      request.headers.set('SIGN', signature);
      request.headers.set('Content-Type', 'application/json');
      request.headers.set('Accept', 'application/json');

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      print('Status Code: ${response.statusCode}');
      print('Response Body: $responseBody');

      if (response.statusCode == 200) {
        final data = jsonDecode(responseBody);
        print('✅ Success! Futures Account:');
        print('  Total: ${data['total']}');
        print('  Available: ${data['available']}');
        print('  Unrealized PnL: ${data['unrealised_pnl']}');
      } else {
        print('❌ Failed: ${response.statusCode}');
      }

      client.close();
    } catch (e) {
      print('❌ Error: $e');
    }
  }

  // Test with different header styles
  static Future<void> testWithDifferentHeaders() async {
    print('\n🔥 Testing with different header styles...');

    final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
    final method = 'GET';
    final path = '/api/v4/spot/accounts';
    final queryString = '';
    final body = '';

    final signature = createSignature(
      method,
      path,
      queryString,
      body,
      timestamp,
    );

    // Test 1: Lowercase headers
    print('\n--- Test 1: Lowercase headers ---');
    await _testWithHeaders({
      'key': apiKey,
      'timestamp': timestamp.toString(),
      'sign': signature,
    }, path);

    // Test 2: Uppercase headers
    print('\n--- Test 2: Uppercase headers ---');
    await _testWithHeaders({
      'KEY': apiKey,
      'TIMESTAMP': timestamp.toString(),
      'SIGN': signature,
    }, path);

    // Test 3: Mixed case headers
    print('\n--- Test 3: Mixed case headers ---');
    await _testWithHeaders({
      'Key': apiKey,
      'Timestamp': timestamp.toString(),
      'Sign': signature,
    }, path);
  }

  static Future<void> _testWithHeaders(
    Map<String, String> headers,
    String path,
  ) async {
    try {
      final client = HttpClient();
      final request = await client.getUrl(Uri.parse('$baseUrl$path'));

      headers.forEach((key, value) {
        request.headers.set(key, value);
      });
      request.headers.set('Content-Type', 'application/json');
      request.headers.set('Accept', 'application/json');

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      print('Status Code: ${response.statusCode}');
      if (response.statusCode != 200) {
        print('Response: $responseBody');
      } else {
        print('✅ Success!');
      }

      client.close();
    } catch (e) {
      print('❌ Error: $e');
    }
  }

  // Test public endpoint to confirm server connectivity
  static Future<void> testPublicEndpoint() async {
    print('\n🔥 Testing Public Endpoint...');

    try {
      final client = HttpClient();
      final request = await client.getUrl(
        Uri.parse('$baseUrl/api/v4/spot/currencies'),
      );

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      print('Status Code: ${response.statusCode}');

      if (response.statusCode == 200) {
        final data = jsonDecode(responseBody) as List;
        print('✅ Public API works! Found ${data.length} currencies');
      } else {
        print('❌ Public API failed: $responseBody');
      }

      client.close();
    } catch (e) {
      print('❌ Error: $e');
    }
  }
}

void main() async {
  print('🚀 Gate.io Go Style API Test');
  print('============================');

  // Test public endpoint first
  await GateGoStyleTest.testPublicEndpoint();

  // Test different header styles
  await GateGoStyleTest.testWithDifferentHeaders();

  // Test specific endpoints
  await GateGoStyleTest.testSpotAccount();
  await GateGoStyleTest.testFuturesAccount();

  print('\n✅ All tests completed');
}
