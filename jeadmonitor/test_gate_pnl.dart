import 'dart:convert';
import 'dart:io';
import 'package:crypto/crypto.dart';

class GateTestClient {
  final String apiKey = '472667d2eebf07b206f8599bf35bbcb1';
  final String apiSecret =
      'e24966020b7241e521966b87deb401d6dac0f90a1ccbf3f85db54ecc09faf33f';
  final String baseUrl = 'https://api.gateio.ws';

  String _createSignature(
    String method,
    String url,
    String queryString,
    String bodyHash,
    int timestamp,
  ) {
    final message = '$method\n$url\n$queryString\n$bodyHash\n$timestamp';
    final key = utf8.encode(apiSecret);
    final messageBytes = utf8.encode(message);
    final hmac = Hmac(sha512, key);
    final digest = hmac.convert(messageBytes);
    return digest.toString();
  }

  Future<void> testFuturesAccount() async {
    try {
      final client = HttpClient();
      final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
      const method = 'GET';
      const url = '/api/v4/futures/usdt/accounts';
      const queryString = '';
      const body = '';
      final bodyHash = sha512.convert(utf8.encode(body)).toString();

      final signature = _createSignature(
        method,
        url,
        queryString,
        bodyHash,
        timestamp,
      );

      final uri = Uri.parse('$baseUrl$url');
      final request = await client.getUrl(uri);
      request.headers.add('KEY', apiKey);
      request.headers.add('Timestamp', timestamp.toString());
      request.headers.add('SIGN', signature);

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      if (response.statusCode == 200) {
        final data = json.decode(responseBody);
        print('=== Gate.io Futures Account ===');
        print('Total: \$${data['total']}');
        print('Unrealized PnL: \$${data['unrealised_pnl']}');
        print('Available: \$${data['available']}');
        print('Position Margin: \$${data['position_margin']}');
        print('Order Margin: \$${data['order_margin']}');
      } else {
        print('Error: ${response.statusCode}');
        print('Response: $responseBody');
      }

      client.close();
    } catch (e) {
      print('Error: $e');
    }
  }

  Future<void> testFuturesPositions() async {
    try {
      final client = HttpClient();
      final timestamp = DateTime.now().millisecondsSinceEpoch ~/ 1000;
      const method = 'GET';
      const url = '/api/v4/futures/usdt/positions';
      const queryString = '';
      const body = '';
      final bodyHash = sha512.convert(utf8.encode(body)).toString();

      final signature = _createSignature(
        method,
        url,
        queryString,
        bodyHash,
        timestamp,
      );

      final uri = Uri.parse('$baseUrl$url');
      final request = await client.getUrl(uri);
      request.headers.add('KEY', apiKey);
      request.headers.add('Timestamp', timestamp.toString());
      request.headers.add('SIGN', signature);

      final response = await request.close();
      final responseBody = await response.transform(utf8.decoder).join();

      if (response.statusCode == 200) {
        final positions = json.decode(responseBody);
        print('\n=== Gate.io Active Positions ===');
        for (var position in positions) {
          final size = double.tryParse(position['size'].toString()) ?? 0.0;
          if (size != 0.0) {
            print(
              '${position['contract']}: ${position['size']} @ \$${position['mark_price']}',
            );
            print('  - Unrealized PnL: \$${position['unrealised_pnl']}');
            print('  - Entry Price: \$${position['entry_price']}');
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
  final client = GateTestClient();
  await client.testFuturesAccount();
  await client.testFuturesPositions();
}
