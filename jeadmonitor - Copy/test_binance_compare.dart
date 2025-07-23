import 'dart:convert';
import 'dart:io';
import 'package:crypto/crypto.dart';

Future<void> main() async {
  const apiKey =
      'xz3oxKAcLPRJxgQvBOSqNGyJipzOoiAvRXwvaHZBYSyHFEfegRwV37Blmwwpxp4z';
  const apiSecret =
      'h9LCIdryqVCEPn6JgGSnmycgtc7FRtKnrSVcqOFVmJ3bOChZDvGBcOCnCdOhvFlT';

  print('🔍 Binance Futures - Current vs Correct Implementation');
  print('=====================================================');

  try {
    final client = BinanceFuturesDebugClient(
      apiKey: apiKey,
      apiSecret: apiSecret,
    );

    // Test current implementation (wrong)
    print('\n❌ CURRENT IMPLEMENTATION (เหมือนใน Flutter app):');
    final currentTotal = await client.getCurrentImplementationTotal();
    print('Current Total: \$${currentTotal.toStringAsFixed(2)}');

    // Test correct implementation
    print('\n✅ CORRECT IMPLEMENTATION:');
    final correctTotal = await client.getCorrectImplementationTotal();
    print('Correct Total: \$${correctTotal.toStringAsFixed(2)}');

    print('\n📊 COMPARISON:');
    print('Current (Wrong): \$${currentTotal.toStringAsFixed(2)}');
    print('Correct (Right): \$${correctTotal.toStringAsFixed(2)}');
    print('Difference: \$${(correctTotal - currentTotal).toStringAsFixed(2)}');
  } catch (e) {
    print('❌ Error: $e');
  }
}

class BinanceFuturesDebugClient {
  static const String futuresUrl = 'https://fapi.binance.com';
  final String apiKey;
  final String apiSecret;

  BinanceFuturesDebugClient({required this.apiKey, required this.apiSecret});

  String _createSignature(String query) {
    var key = utf8.encode(apiSecret);
    var bytes = utf8.encode(query);
    var hmacSha256 = Hmac(sha256, key);
    var digest = hmacSha256.convert(bytes);
    return digest.toString();
  }

  Future<Map<String, dynamic>> getFuturesAccount() async {
    final timestamp = DateTime.now().millisecondsSinceEpoch;
    final query = 'timestamp=$timestamp&recvWindow=5000';
    final signature = _createSignature(query);

    final url = '$futuresUrl/fapi/v2/account?$query&signature=$signature';

    final request = await HttpClient().getUrl(Uri.parse(url));
    request.headers.set('X-MBX-APIKEY', apiKey);

    final response = await request.close();
    final responseBody = await response.transform(utf8.decoder).join();

    if (response.statusCode == 200) {
      return jsonDecode(responseBody);
    } else {
      throw Exception(
        'Failed to get futures account: ${response.statusCode} - $responseBody',
      );
    }
  }

  // Current implementation (wrong) - เหมือนใน Flutter app
  Future<double> getCurrentImplementationTotal() async {
    print('📱 Simulating Flutter app logic...');

    final accountInfo = await getFuturesAccount();
    double total = 0.0;

    // ใช้ assets array และ mock price (ผิด)
    for (var balance in accountInfo['assets'] ?? []) {
      final double walletBalance =
          double.tryParse(balance['walletBalance']?.toString() ?? '0') ?? 0.0;

      if (walletBalance > 0) {
        final mockPrice = 50.0; // Mock price ใน Flutter app
        final usdValue = walletBalance * mockPrice;
        total += usdValue;

        print(
          '  ${balance['asset']}: ${walletBalance} × \$${mockPrice} = \$${usdValue.toStringAsFixed(2)}',
        );
      }
    }

    return total;
  }

  // Correct implementation
  Future<double> getCorrectImplementationTotal() async {
    print('✅ Using correct totalWalletBalance...');

    final accountInfo = await getFuturesAccount();

    // ใช้ totalWalletBalance โดยตรง (ถูก)
    final totalWalletBalance =
        double.tryParse(accountInfo['totalWalletBalance']?.toString() ?? '0') ??
        0.0;
    final totalUnrealizedProfit =
        double.tryParse(
          accountInfo['totalUnrealizedProfit']?.toString() ?? '0',
        ) ??
        0.0;
    final availableBalance =
        double.tryParse(accountInfo['availableBalance']?.toString() ?? '0') ??
        0.0;

    print('  Total Wallet Balance: \$${totalWalletBalance}');
    print('  Unrealized Profit: \$${totalUnrealizedProfit}');
    print('  Available Balance: \$${availableBalance}');

    return totalWalletBalance;
  }
}
