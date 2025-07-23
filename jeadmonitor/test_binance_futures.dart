import 'dart:convert';
import 'dart:io';
import 'package:crypto/crypto.dart';

Future<void> main() async {
  const apiKey =
      'xz3oxKAcLPRJxgQvBOSqNGyJipzOoiAvRXwvaHZBYSyHFEfegRwV37Blmwwpxp4z';
  const apiSecret =
      'h9LCIdryqVCEPn6JgGSnmycgtc7FRtKnrSVcqOFVmJ3bOChZDvGBcOCnCdOhvFlT';

  print('🚀 Binance Futures API Test');
  print('============================');

  try {
    final client = BinanceFuturesClient(apiKey: apiKey, apiSecret: apiSecret);

    print('\n📊 Testing Futures Account Balance...');
    final account = await client.getFuturesAccount();
    print('Account Balance: ${account['totalWalletBalance']} USDT');
    print('Unrealized PnL: ${account['totalUnrealizedProfit']} USDT');
    print('Available Balance: ${account['availableBalance']} USDT');

    print('\n📈 Testing Futures Positions...');
    final positions = await client.getFuturesPositions();
    print('Total positions: ${positions.length}');

    double totalValue = 0;
    double totalPnL = 0;
    int activePositions = 0;

    for (var position in positions) {
      final size =
          double.tryParse(position['positionAmt']?.toString() ?? '0') ?? 0;
      final price =
          double.tryParse(position['markPrice']?.toString() ?? '0') ?? 0;
      final pnl =
          double.tryParse(position['unRealizedProfit']?.toString() ?? '0') ?? 0;

      if (size.abs() > 0.0001) {
        activePositions++;
        final value = size.abs() * price;
        totalValue += value;
        totalPnL += pnl;

        print(
          '${position['symbol']}: ${size > 0 ? 'LONG' : 'SHORT'} ${size.abs()} @ \$${price.toStringAsFixed(4)}',
        );
        print(
          '  Value: \$${value.toStringAsFixed(2)}, PnL: \$${pnl.toStringAsFixed(2)}',
        );
      }
    }

    print('\n💰 Summary:');
    print('Total Wallet Balance: \$${account['totalWalletBalance']}');
    print('Total Position Value: \$${totalValue.toStringAsFixed(2)}');
    print('Total Unrealized PnL: \$${totalPnL.toStringAsFixed(2)}');
    print('Active Positions: $activePositions');
  } catch (e) {
    print('❌ Error: $e');
  }
}

class BinanceFuturesClient {
  static const String baseUrl = 'https://fapi.binance.com';
  final String apiKey;
  final String apiSecret;

  BinanceFuturesClient({required this.apiKey, required this.apiSecret});

  String _createSignature(String query) {
    var key = utf8.encode(apiSecret);
    var bytes = utf8.encode(query);
    var hmacSha256 = Hmac(sha256, key);
    var digest = hmacSha256.convert(bytes);
    return digest.toString();
  }

  Future<Map<String, dynamic>> getFuturesAccount() async {
    final timestamp = DateTime.now().millisecondsSinceEpoch;
    final query = 'timestamp=$timestamp';
    final signature = _createSignature(query);

    final url = '$baseUrl/fapi/v2/account?$query&signature=$signature';

    print('🔍 Request URL: $url');

    final request = await HttpClient().getUrl(Uri.parse(url));
    request.headers.set('X-MBX-APIKEY', apiKey);

    final response = await request.close();
    final responseBody = await response.transform(utf8.decoder).join();

    print('📡 Response Status: ${response.statusCode}');
    print(
      '📄 Response Body: ${responseBody.substring(0, responseBody.length > 500 ? 500 : responseBody.length)}...',
    );

    if (response.statusCode == 200) {
      return jsonDecode(responseBody);
    } else {
      throw Exception(
        'Failed to get futures account: ${response.statusCode} - $responseBody',
      );
    }
  }

  Future<List<dynamic>> getFuturesPositions() async {
    final timestamp = DateTime.now().millisecondsSinceEpoch;
    final query = 'timestamp=$timestamp';
    final signature = _createSignature(query);

    final url = '$baseUrl/fapi/v2/positionRisk?$query&signature=$signature';

    final request = await HttpClient().getUrl(Uri.parse(url));
    request.headers.set('X-MBX-APIKEY', apiKey);

    final response = await request.close();
    final responseBody = await response.transform(utf8.decoder).join();

    print('📡 Positions Response Status: ${response.statusCode}');

    if (response.statusCode == 200) {
      return jsonDecode(responseBody);
    } else {
      throw Exception(
        'Failed to get futures positions: ${response.statusCode} - $responseBody',
      );
    }
  }
}
