import 'dart:convert';
import 'dart:io';
import 'package:crypto/crypto.dart';

class BinanceTestClient {
  final String apiKey =
      'xz3oxKAcLPRJxgQvBOSqNGyJipzOoiAvRXwvaHZBYSyHFEfegRwV37Blmwwpxp4z';
  final String apiSecret =
      'h9LCIdryqVCEPn6JgGSnmycgtc7FRtKnrSVcqOFVmJ3bOChZDvGBcOCnCdOhvFlT';
  final String baseUrl = 'https://api.binance.com';
  final String futuresUrl = 'https://fapi.binance.com';

  String _createSignature(String query) {
    final key = utf8.encode(apiSecret);
    final message = utf8.encode(query);
    final hmac = Hmac(sha256, key);
    final digest = hmac.convert(message);
    return digest.toString();
  }

  Future<void> testFuturesAccount() async {
    try {
      final client = HttpClient();
      final timestamp = DateTime.now().millisecondsSinceEpoch.toString();
      final query = 'timestamp=$timestamp&recvWindow=5000';
      final signature = _createSignature(query);

      final uri = Uri.parse(
        '$futuresUrl/fapi/v2/account?$query&signature=$signature',
      );
      final request = await client.getUrl(uri);
      request.headers.add('X-MBX-APIKEY', apiKey);
      request.headers.add('Content-Type', 'application/json');

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      if (response.statusCode == 200) {
        final data = json.decode(responseBody);
        print('=== Binance Futures Account Info ===');
        print('Total Wallet Balance: \$${data['totalWalletBalance']}');
        print('Total Unrealized Profit: \$${data['totalUnrealizedProfit']}');
        print('Total Margin Balance: \$${data['totalMarginBalance']}');
        print('Available Balance: \$${data['availableBalance']}');
        print(
          'Total Position Initial Margin: \$${data['totalPositionInitialMargin']}',
        );

        // Show individual positions with PnL
        if (data['positions'] != null) {
          print('\n=== Active Positions ===');
          for (var position in data['positions']) {
            final positionAmt =
                double.tryParse(position['positionAmt'].toString()) ?? 0.0;
            if (positionAmt != 0.0) {
              print(
                '${position['symbol']}: ${position['positionAmt']} @ \$${position['markPrice']}',
              );
              print('  - Unrealized PnL: \$${position['unrealizedProfit']}');
              print('  - Entry Price: \$${position['entryPrice']}');
            }
          }
        }
      } else {
        print('Error: ${response.statusCode}');
        print('Response: $responseBody');
      }

      client.close();
    } catch (e) {
      print('Error: $e');
    }
  }
}

void main() async {
  final client = BinanceTestClient();
  await client.testFuturesAccount();
}
