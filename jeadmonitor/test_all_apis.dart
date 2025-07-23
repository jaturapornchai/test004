import 'dart:io';
import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:http/http.dart' as http;

class BinanceAPITester {
  final String apiKey;
  final String apiSecret;
  final String baseUrl = 'https://api.binance.com';
  final String futuresUrl = 'https://fapi.binance.com';

  BinanceAPITester(this.apiKey, this.apiSecret);

  String _generateSignature(String query) {
    var key = utf8.encode(apiSecret);
    var bytes = utf8.encode(query);
    var hmacSha256 = Hmac(sha256, key);
    return hmacSha256.convert(bytes).toString();
  }

  Future<void> testSpotAccount() async {
    print('\n=== Testing Binance Spot Account ===');
    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch;
      final query = 'timestamp=$timestamp';
      final signature = _generateSignature(query);

      final url = '$baseUrl/api/v3/account?$query&signature=$signature';

      final response = await http.get(
        Uri.parse(url),
        headers: {'X-MBX-APIKEY': apiKey, 'Content-Type': 'application/json'},
      );

      print('Status Code: ${response.statusCode}');
      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        final balances = (data['balances'] as List)
            .where(
              (balance) =>
                  double.parse(balance['free']) > 0 ||
                  double.parse(balance['locked']) > 0,
            )
            .toList();

        print('✅ Binance Spot API working!');
        print('Account Type: ${data['accountType']}');
        print('Assets with balance: ${balances.length}');

        for (var balance in balances.take(5)) {
          final free = double.parse(balance['free']);
          final locked = double.parse(balance['locked']);
          if (free > 0 || locked > 0) {
            print('  ${balance['asset']}: Free: $free, Locked: $locked');
          }
        }
      } else {
        print('❌ Error: ${response.body}');
      }
    } catch (e) {
      print('❌ Exception: $e');
    }
  }

  Future<void> testFuturesAccount() async {
    print('\n=== Testing Binance Futures Account ===');
    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch;
      final query = 'timestamp=$timestamp';
      final signature = _generateSignature(query);

      final url = '$futuresUrl/fapi/v2/account?$query&signature=$signature';

      final response = await http.get(
        Uri.parse(url),
        headers: {'X-MBX-APIKEY': apiKey, 'Content-Type': 'application/json'},
      );

      print('Status Code: ${response.statusCode}');
      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        final assets = (data['assets'] as List)
            .where((asset) => double.parse(asset['walletBalance']) > 0)
            .toList();

        print('✅ Binance Futures API working!');
        print('Total Wallet Balance: ${data['totalWalletBalance']} USDT');
        print('Total Unrealized PNL: ${data['totalUnrealizedProfit']} USDT');
        print('Assets with balance: ${assets.length}');

        for (var asset in assets.take(5)) {
          final balance = double.parse(asset['walletBalance']);
          final pnl = double.parse(asset['unrealizedProfit']);
          print('  ${asset['asset']}: Balance: $balance, PNL: $pnl');
        }

        final positions = (data['positions'] as List)
            .where((pos) => double.parse(pos['positionAmt']) != 0)
            .toList();

        if (positions.isNotEmpty) {
          print('Open Positions: ${positions.length}');
          for (var pos in positions.take(3)) {
            print(
              '  ${pos['symbol']}: ${pos['positionAmt']} @ ${pos['entryPrice']}',
            );
          }
        }
      } else {
        print('❌ Error: ${response.body}');
      }
    } catch (e) {
      print('❌ Exception: $e');
    }
  }
}

class GateAPITester {
  final String apiKey;
  final String apiSecret;
  final String baseUrl = 'https://api.gateio.ws';

  GateAPITester(this.apiKey, this.apiSecret);

  String _generateSignature(
    String method,
    String url,
    String queryString,
    String body,
  ) {
    final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
    final hashedBody = sha512.convert(utf8.encode(body)).toString();
    final signString = '$method\n$url\n$queryString\n$hashedBody\n$timestamp';

    final key = utf8.encode(apiSecret);
    final bytes = utf8.encode(signString);
    final hmacSha512 = Hmac(sha512, key);
    return hmacSha512.convert(bytes).toString();
  }

  Future<void> testSpotAccount() async {
    print('\n=== Testing Gate.io Spot Account ===');
    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
      final method = 'GET';
      final url = '/api/v4/spot/accounts';
      final queryString = '';
      final body = '';

      final signature = _generateSignature(method, url, queryString, body);

      final response = await http.get(
        Uri.parse('$baseUrl$url'),
        headers: {
          'KEY': apiKey,
          'SIGN': signature,
          'Timestamp': timestamp.toString(),
          'Content-Type': 'application/json',
        },
      );

      print('Status Code: ${response.statusCode}');
      if (response.statusCode == 200) {
        final data = json.decode(response.body) as List;
        final balances = data
            .where(
              (balance) =>
                  double.parse(balance['available']) > 0 ||
                  double.parse(balance['locked']) > 0,
            )
            .toList();

        print('✅ Gate.io Spot API working!');
        print('Assets with balance: ${balances.length}');

        for (var balance in balances.take(5)) {
          final available = double.parse(balance['available']);
          final locked = double.parse(balance['locked']);
          if (available > 0 || locked > 0) {
            print(
              '  ${balance['currency']}: Available: $available, Locked: $locked',
            );
          }
        }
      } else {
        print('❌ Error: ${response.body}');
        print('Headers sent:');
        print('  KEY: $apiKey');
        print('  Timestamp: $timestamp');
        print('  Signature: $signature');
      }
    } catch (e) {
      print('❌ Exception: $e');
    }
  }

  Future<void> testFuturesAccount() async {
    print('\n=== Testing Gate.io Futures Account ===');
    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
      final method = 'GET';
      final url = '/api/v4/futures/usdt/accounts';
      final queryString = '';
      final body = '';

      final signature = _generateSignature(method, url, queryString, body);

      final response = await http.get(
        Uri.parse('$baseUrl$url'),
        headers: {
          'KEY': apiKey,
          'SIGN': signature,
          'Timestamp': timestamp.toString(),
          'Content-Type': 'application/json',
        },
      );

      print('Status Code: ${response.statusCode}');
      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        print('✅ Gate.io Futures API working!');
        print('Total: ${data['total']} USDT');
        print('Available: ${data['available']} USDT');
        print('Unrealized PNL: ${data['unrealised_pnl']} USDT');

        // Test positions
        await _testFuturesPositions();
      } else {
        print('❌ Error: ${response.body}');
      }
    } catch (e) {
      print('❌ Exception: $e');
    }
  }

  Future<void> _testFuturesPositions() async {
    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
      final method = 'GET';
      final url = '/api/v4/futures/usdt/positions';
      final queryString = '';
      final body = '';

      final signature = _generateSignature(method, url, queryString, body);

      final response = await http.get(
        Uri.parse('$baseUrl$url'),
        headers: {
          'KEY': apiKey,
          'SIGN': signature,
          'Timestamp': timestamp.toString(),
          'Content-Type': 'application/json',
        },
      );

      if (response.statusCode == 200) {
        final positions = json.decode(response.body) as List;
        final openPositions = positions
            .where((pos) => int.parse(pos['size']) != 0)
            .toList();

        print('Open Futures Positions: ${openPositions.length}');
        for (var pos in openPositions.take(3)) {
          print(
            '  ${pos['contract']}: Size: ${pos['size']}, Value: ${pos['value']}',
          );
        }
      }
    } catch (e) {
      print('❌ Exception getting positions: $e');
    }
  }
}

void main() async {
  // Load environment variables
  final envFile = File('.env');
  final envContent = await envFile.readAsString();
  final envLines = envContent
      .split('\n')
      .where((line) => line.contains('='))
      .toList();

  String? gateApiKey, gateApiSecret, binanceApiKey, binanceApiSecret;

  for (final line in envLines) {
    final parts = line.split('=');
    if (parts.length == 2) {
      final key = parts[0].trim();
      final value = parts[1].trim();

      switch (key) {
        case 'GATE_API_KEY':
          gateApiKey = value;
          break;
        case 'GATE_API_SECRET':
          gateApiSecret = value;
          break;
        case 'BINANCE_API_KEY':
          binanceApiKey = value;
          break;
        case 'BINANCE_API_SECRET':
          binanceApiSecret = value;
          break;
      }
    }
  }

  print('🚀 Testing All APIs - Spot & Futures Markets');
  print('============================================');

  if (binanceApiKey != null && binanceApiSecret != null) {
    final binanceTester = BinanceAPITester(binanceApiKey, binanceApiSecret);
    await binanceTester.testSpotAccount();
    await binanceTester.testFuturesAccount();
  } else {
    print('❌ Binance API keys not found');
  }

  if (gateApiKey != null && gateApiSecret != null) {
    final gateTester = GateAPITester(gateApiKey, gateApiSecret);
    await gateTester.testSpotAccount();
    await gateTester.testFuturesAccount();
  } else {
    print('❌ Gate.io API keys not found');
  }

  print('\n✅ API Testing Complete!');
}
