import 'dart:convert';
import 'dart:io';
import 'package:crypto/crypto.dart';

class GateCorrectStyleTest {
  static const String baseUrl = 'https://api.gateio.ws';
  static const String apiKey = '472667d2eebf07b206f8599bf35bbcb1';
  static const String apiSecret =
      'e24966020b7241e521966b87deb401d6dac0f90a1ccbf3f85db54ecc09faf33f';

  // สร้าง signature ตามรูปแบบ Go client ที่ถูกต้อง
  static String createSignature(
    String method,
    String path,
    String queryString,
    String body,
    int timestamp,
  ) {
    // 1. สร้าง SHA512 hash ของ body
    var bodyBytes = utf8.encode(body);
    var bodyHash = sha512.convert(bodyBytes);
    String hashedPayload = bodyHash.toString();

    // 2. สร้าง message string
    String message = '$method\n$path\n$queryString\n$hashedPayload\n$timestamp';

    print('=== Correct Authentication Debug ===');
    print('Method: $method');
    print('Path: $path');
    print('Query String: $queryString');
    print('Body: $body');
    print('Body Hash (SHA512): $hashedPayload');
    print('Timestamp: $timestamp');
    print('Message: "$message"');

    // 3. สร้าง HMAC-SHA512 signature
    var secretBytes = utf8.encode(apiSecret);
    var messageBytes = utf8.encode(message);
    var hmac = Hmac(sha512, secretBytes);
    var signature = hmac.convert(messageBytes);

    print('Signature: $signature');
    print('==============================');

    return signature.toString();
  }

  // Test Spot Account ด้วยวิธีที่ถูกต้อง
  static Future<void> testSpotAccount() async {
    print('\n🔥 Testing Spot Account (Correct Style)...');

    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
      final method = 'GET';
      final path = '/api/v4/spot/accounts';
      final queryString = '';
      final body = ''; // Empty body for GET request

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
      request.headers.set('SIGN', signature);
      request.headers.set('Timestamp', timestamp.toString());
      request.headers.set('Content-Type', 'application/json');
      request.headers.set('Accept', 'application/json');

      print('Request Headers:');
      print('  KEY: $apiKey');
      print('  SIGN: $signature');
      print('  Timestamp: $timestamp');

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      print('Status Code: ${response.statusCode}');

      if (response.statusCode == 200) {
        final data = jsonDecode(responseBody) as List;
        print('✅ SUCCESS! Found ${data.length} accounts');

        double totalValue = 0;
        for (var account in data) {
          final currency = account['currency'];
          final available = double.parse(account['available']);
          final locked = double.parse(account['locked']);
          final total = available + locked;

          if (total > 0) {
            print(
              '  $currency: Available=$available, Locked=$locked, Total=$total',
            );
            totalValue += total; // This is just a sum, not USD value
          }
        }
        print(
          '  Total Balance Items with Value: ${totalValue.toStringAsFixed(8)}',
        );
      } else {
        print('❌ Failed: ${response.statusCode}');
        print('Response: $responseBody');
      }

      client.close();
    } catch (e) {
      print('❌ Error: $e');
    }
  }

  // Test Futures Account ด้วยวิธีที่ถูกต้อง
  static Future<void> testFuturesAccount() async {
    print('\n🔥 Testing Futures Account (Correct Style)...');

    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
      final method = 'GET';
      final path = '/api/v4/futures/usdt/accounts';
      final queryString = '';
      final body = ''; // Empty body for GET request

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
      request.headers.set('SIGN', signature);
      request.headers.set('Timestamp', timestamp.toString());
      request.headers.set('Content-Type', 'application/json');
      request.headers.set('Accept', 'application/json');

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      print('Status Code: ${response.statusCode}');

      if (response.statusCode == 200) {
        final data = jsonDecode(responseBody);
        print('✅ SUCCESS! Futures Account:');
        print('  Total: ${data['total']}');
        print('  Available: ${data['available']}');
        print('  Unrealized PnL: ${data['unrealised_pnl']}');
        print('  Position Margin: ${data['position_margin']}');
        print('  Order Margin: ${data['order_margin']}');
      } else {
        print('❌ Failed: ${response.statusCode}');
        print('Response: $responseBody');
      }

      client.close();
    } catch (e) {
      print('❌ Error: $e');
    }
  }

  // Test Futures Positions
  static Future<void> testFuturesPositions() async {
    print('\n🔥 Testing Futures Positions (Correct Style)...');

    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
      final method = 'GET';
      final path = '/api/v4/futures/usdt/positions';
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
      request.headers.set('SIGN', signature);
      request.headers.set('Timestamp', timestamp.toString());
      request.headers.set('Content-Type', 'application/json');
      request.headers.set('Accept', 'application/json');

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      print('Status Code: ${response.statusCode}');

      if (response.statusCode == 200) {
        final data = jsonDecode(responseBody) as List;
        print('✅ SUCCESS! Found ${data.length} positions');

        for (var position in data) {
          final contract = position['contract'];
          final size = position['size'];
          final value = position['value'];
          final unrealizedPnl = position['unrealised_pnl'];

          if (double.parse(size) != 0) {
            print('  $contract: Size=$size, Value=$value, PnL=$unrealizedPnl');
          }
        }
      } else {
        print('❌ Failed: ${response.statusCode}');
        print('Response: $responseBody');
      }

      client.close();
    } catch (e) {
      print('❌ Error: $e');
    }
  }

  // Test public endpoint เพื่อยืนยันการเชื่อมต่อ
  static Future<void> testPublicEndpoint() async {
    print('\n🔥 Testing Public Endpoint...');

    try {
      final client = HttpClient();
      final request = await client.getUrl(
        Uri.parse('$baseUrl/api/v4/spot/time'),
      );

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      print('Status Code: ${response.statusCode}');

      if (response.statusCode == 200) {
        final data = jsonDecode(responseBody);
        print('✅ Server time: ${data['server_time']}');
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
  print('🚀 Gate.io Correct Authentication Test');
  print('======================================');

  // Test public endpoint first
  await GateCorrectStyleTest.testPublicEndpoint();

  // Test authenticated endpoints with correct signature
  await GateCorrectStyleTest.testSpotAccount();
  await GateCorrectStyleTest.testFuturesAccount();
  await GateCorrectStyleTest.testFuturesPositions();

  print('\n✅ All tests completed');
}
