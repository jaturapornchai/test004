import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:http/http.dart' as http;

// Simple test to check if API keys are working
class ApiTester {
  // Test Binance API
  static Future<void> testBinanceApi() async {
    const apiKey =
        'xz3oxKAcLPRJxgQvBOSqNGyJipzOoiAvRXwvaHZBYSyHFEfegRwV37Blmwwpxp4z';
    const apiSecret =
        'h9LCIdryqVCEPn6JgGSnmycgtc7FRtKnrSVcqOFVmJ3bOChZDvGBcOCnCdOhvFlT';

    try {
      // Test simple ping first
      final pingResponse = await http.get(
        Uri.parse('https://api.binance.com/api/v3/ping'),
      );
      print('Binance Ping: ${pingResponse.statusCode}');

      // Test server time
      final timeResponse = await http.get(
        Uri.parse('https://api.binance.com/api/v3/time'),
      );
      if (timeResponse.statusCode == 200) {
        final timeData = json.decode(timeResponse.body);
        print('Binance Server Time: ${timeData['serverTime']}');

        // Test account info with proper signature
        final timestamp = timeData['serverTime'].toString();
        final query = 'timestamp=$timestamp&recvWindow=5000';

        var key = utf8.encode(apiSecret);
        var bytes = utf8.encode(query);
        var hmacSha256 = Hmac(sha256, key);
        var signature = hmacSha256.convert(bytes).toString();

        final accountResponse = await http.get(
          Uri.parse(
            'https://api.binance.com/api/v3/account?$query&signature=$signature',
          ),
          headers: {'X-MBX-APIKEY': apiKey},
        );

        print(
          'Binance Account: ${accountResponse.statusCode} - ${accountResponse.body.substring(0, 200)}...',
        );
      }
    } catch (e) {
      print('Binance Test Error: $e');
    }
  }

  // Test Gate.io API
  static Future<void> testGateApi() async {
    const apiKey = '472667d2eebf07b206f8599bf35bbcb1';
    const apiSecret =
        'e24966020b7241e521966b87deb401d6dac0f90a1ccbf3f85db54ecc09faf33f';

    try {
      // Test public endpoint first
      final tickerResponse = await http.get(
        Uri.parse('https://api.gateio.ws/api/v4/spot/tickers'),
      );
      print('Gate.io Public API: ${tickerResponse.statusCode}');

      // Test private endpoint with proper hash calculation
      const method = 'GET';
      const path = '/spot/accounts';
      const queryString = '';
      const bodyContent = '';
      final timestamp = (DateTime.now().millisecondsSinceEpoch / 1000)
          .round()
          .toString();

      // Calculate SHA256 hash of body
      var bodyBytes = utf8.encode(bodyContent);
      var bodyHash = sha256.convert(bodyBytes).toString();

      // Gate.io signature format: METHOD\n/PATH\nQUERY_STRING\nBODY_HASH\nTIMESTAMP
      final payloadToSign =
          '$method\n$path\n$queryString\n$bodyHash\n$timestamp';

      print('Gate.io Body Hash: $bodyHash');
      print('Gate.io Payload to sign: "$payloadToSign"');

      var key = utf8.encode(apiSecret);
      var bytes = utf8.encode(payloadToSign);
      var hmacSha512 = Hmac(sha512, key);
      var signature = hmacSha512.convert(bytes).toString();

      print('Gate.io Timestamp: $timestamp');
      print('Gate.io Signature: $signature');

      final accountResponse = await http.get(
        Uri.parse('https://api.gateio.ws/api/v4$path'),
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'KEY': apiKey,
          'Timestamp': timestamp,
          'SIGN': signature,
        },
      );

      print(
        'Gate.io Account: ${accountResponse.statusCode} - ${accountResponse.body}',
      );

      // If still 401, let's try a different approach - test with balances endpoint
      if (accountResponse.statusCode == 401) {
        print('\nTrying different endpoint: /wallet/total_balance');
        const balancePath = '/wallet/total_balance';
        final balancePayload =
            '$method\n$balancePath\n$queryString\n$bodyHash\n$timestamp';

        var balanceBytes = utf8.encode(balancePayload);
        var balanceSignature = Hmac(
          sha512,
          key,
        ).convert(balanceBytes).toString();

        final balanceResponse = await http.get(
          Uri.parse('https://api.gateio.ws/api/v4$balancePath'),
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'KEY': apiKey,
            'Timestamp': timestamp,
            'SIGN': balanceSignature,
          },
        );

        print(
          'Gate.io Balance: ${balanceResponse.statusCode} - ${balanceResponse.body}',
        );
      }
    } catch (e) {
      print('Gate.io Test Error: $e');
    }
  }
}

void main() async {
  print('Testing API Keys...');

  print('\n=== Testing Binance ===');
  await ApiTester.testBinanceApi();

  print('\n=== Testing Gate.io ===');
  await ApiTester.testGateApi();
}
