import 'dart:convert';
import 'dart:io';
import 'package:crypto/crypto.dart';
import 'package:http/http.dart' as http;

void main() async {
  print('=== Gate.io API Signature Test ===\n');

  const apiKey = '472667d2eebf07b206f8599bf35bbcb1';
  const apiSecret =
      'e24966020b7241e521966b87deb401d6dac0f90a1ccbf3f85db54ecc09faf33f';

  await testGateSpotSignature(apiKey, apiSecret);
  await testGateFuturesSignature(apiKey, apiSecret);

  exit(0);
}

Future<void> testGateSpotSignature(String apiKey, String apiSecret) async {
  print('🔑 Testing Gate.io Spot Signature...\n');

  try {
    final timestamp = (DateTime.now().millisecondsSinceEpoch / 1000)
        .round()
        .toString();
    const method = 'GET';
    const path = '/api/v4/spot/accounts';
    const queryString = '';
    const body = '';

    // Create signature string
    final payloadHash = sha512.convert(utf8.encode(body)).toString();
    final signString = '$method\n$path\n$queryString\n$payloadHash\n$timestamp';

    print('📝 Signature components:');
    print('Method: $method');
    print('Path: $path');
    print('Query: $queryString');
    print('Body: $body');
    print('Body Hash: $payloadHash');
    print('Timestamp: $timestamp');
    print('Sign String: $signString');

    final signature = Hmac(
      sha512,
      utf8.encode(apiSecret),
    ).convert(utf8.encode(signString)).toString();

    print('Signature: $signature\n');

    final headers = {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'KEY': apiKey,
      'Timestamp': timestamp,
      'SIGN': signature,
    };

    print('📤 Headers:');
    headers.forEach((key, value) => print('$key: $value'));

    final url = 'https://api.gateio.ws$path';
    print('\n🌐 Making request to: $url');

    final response = await http.get(Uri.parse(url), headers: headers);

    print('\n📥 Response:');
    print('Status: ${response.statusCode}');
    print('Body: ${response.body}');

    if (response.statusCode == 200) {
      print('✅ Gate.io Spot API Success!');
      final data = jsonDecode(response.body) as List;
      print('📊 Found ${data.length} spot accounts');

      double totalValue = 0.0;
      for (var account in data) {
        final available =
            double.tryParse(account['available'].toString()) ?? 0.0;
        final locked = double.tryParse(account['locked'].toString()) ?? 0.0;
        final total = available + locked;

        if (total > 0.0001) {
          print('  ${account['currency']}: ${total.toStringAsFixed(8)}');
          if (account['currency'] == 'USDT') totalValue += total;
        }
      }
      print('💰 Estimated USDT value: \$${totalValue.toStringAsFixed(2)}');
    } else {
      print('❌ Gate.io Spot API Failed!');
    }
  } catch (e) {
    print('❌ Error: $e');
  }
}

Future<void> testGateFuturesSignature(String apiKey, String apiSecret) async {
  print('\n🔑 Testing Gate.io Futures Signature...\n');

  try {
    final timestamp = (DateTime.now().millisecondsSinceEpoch / 1000)
        .round()
        .toString();
    const method = 'GET';
    const path = '/api/v4/futures/usdt/accounts';
    const queryString = '';
    const body = '';

    // Create signature string
    final payloadHash = sha512.convert(utf8.encode(body)).toString();
    final signString = '$method\n$path\n$queryString\n$payloadHash\n$timestamp';

    print('📝 Signature components:');
    print('Method: $method');
    print('Path: $path');
    print('Query: $queryString');
    print('Body: $body');
    print('Body Hash: $payloadHash');
    print('Timestamp: $timestamp');
    print('Sign String: $signString');

    final signature = Hmac(
      sha512,
      utf8.encode(apiSecret),
    ).convert(utf8.encode(signString)).toString();

    print('Signature: $signature\n');

    final headers = {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'KEY': apiKey,
      'Timestamp': timestamp,
      'SIGN': signature,
    };

    print('📤 Headers:');
    headers.forEach((key, value) => print('$key: $value'));

    final url = 'https://api.gateio.ws$path';
    print('\n🌐 Making request to: $url');

    final response = await http.get(Uri.parse(url), headers: headers);

    print('\n📥 Response:');
    print('Status: ${response.statusCode}');
    print('Body: ${response.body}');

    if (response.statusCode == 200) {
      print('✅ Gate.io Futures API Success!');
      final data = jsonDecode(response.body);
      print('📊 Futures account data: $data');

      if (data['total'] != null) {
        final total = double.tryParse(data['total'].toString()) ?? 0.0;
        print('💰 Total futures balance: \$${total.toStringAsFixed(2)}');
      }
    } else {
      print('❌ Gate.io Futures API Failed!');
    }
  } catch (e) {
    print('❌ Error: $e');
  }
}
