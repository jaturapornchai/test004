import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import 'lib/services/binance_api_client.dart';
import 'lib/services/gate_api_client.dart';

void main() async {
  print('=== Real Portfolio Balance Calculator ===\n');

  final binanceClient = BinanceApiClient(
    apiKey: 'xz3oxKAcLPRJxgQvBOSqNGyJipzOoiAvRXwvaHZBYSyHFEfegRwV37Blmwwpxp4z',
    apiSecret:
        'h9LCIdryqVCEPn6JgGSnmycgtc7FRtKnrSVcqOFVmJ3bOChZDvGBcOCnCdOhvFlT',
  );

  final gateClient = GateApiClient(
    apiKey: '472667d2eebf07b206f8599bf35bbcb1',
    apiSecret:
        'e24966020b7241e521966b87deb401d6dac0f90a1ccbf3f85db54ecc09faf33f',
  );

  // Get current prices from CoinGecko
  final prices = await getCryptoPrices();
  print('📈 Current crypto prices loaded\n');

  double totalPortfolioValue = 0.0;

  // Test Binance
  print('🟡 Binance Portfolio:');
  final binanceSpotValue = await calculateBinanceSpot(binanceClient, prices);
  final binanceFuturesValue = await calculateBinanceFutures(binanceClient);
  final binanceTotal = binanceSpotValue + binanceFuturesValue;

  print('  💰 Spot: \$${binanceSpotValue.toStringAsFixed(2)}');
  print('  💰 Futures: \$${binanceFuturesValue.toStringAsFixed(2)}');
  print('  📊 Binance Total: \$${binanceTotal.toStringAsFixed(2)}\n');

  totalPortfolioValue += binanceTotal;

  // Test Gate.io
  print('🟢 Gate.io Portfolio:');
  final gateSpotValue = await calculateGateSpot(gateClient, prices);
  final gateFuturesValue = await calculateGateFutures(gateClient);
  final gateTotal = gateSpotValue + gateFuturesValue;

  print('  💰 Spot: \$${gateSpotValue.toStringAsFixed(2)}');
  print('  💰 Futures: \$${gateFuturesValue.toStringAsFixed(2)}');
  print('  📊 Gate.io Total: \$${gateTotal.toStringAsFixed(2)}\n');

  totalPortfolioValue += gateTotal;

  print('=' * 50);
  print(
    '🎯 TOTAL PORTFOLIO VALUE: \$${totalPortfolioValue.toStringAsFixed(2)}',
  );
  print('=' * 50);

  exit(0);
}

Future<Map<String, double>> getCryptoPrices() async {
  try {
    const url =
        'https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum,cardano,binancecoin,usd-coin,tether,solana&vs_currencies=usd';
    final response = await http.get(Uri.parse(url));

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body) as Map<String, dynamic>;
      return {
        'BTC': (data['bitcoin']?['usd'] ?? 45000).toDouble(),
        'ETH': (data['ethereum']?['usd'] ?? 2500).toDouble(),
        'ADA': (data['cardano']?['usd'] ?? 0.4).toDouble(),
        'BNB': (data['binancecoin']?['usd'] ?? 300).toDouble(),
        'USDC': (data['usd-coin']?['usd'] ?? 1.0).toDouble(),
        'USDT': (data['tether']?['usd'] ?? 1.0).toDouble(),
        'SOL': (data['solana']?['usd'] ?? 100).toDouble(),
      };
    }
  } catch (e) {
    print('⚠️ Failed to fetch real prices, using defaults: $e');
  }

  // Fallback prices
  return {
    'BTC': 45000.0,
    'ETH': 2500.0,
    'ADA': 0.4,
    'BNB': 300.0,
    'USDC': 1.0,
    'USDT': 1.0,
    'SOL': 100.0,
  };
}

Future<double> calculateBinanceSpot(
  BinanceApiClient client,
  Map<String, double> prices,
) async {
  try {
    final account = await client.getAccountInfo();
    double totalValue = 0.0;

    print('  📋 Binance Spot Assets:');
    for (var balance in account['balances']) {
      final symbol = balance['asset'] as String;
      final double free = double.tryParse(balance['free'].toString()) ?? 0.0;
      final double locked =
          double.tryParse(balance['locked'].toString()) ?? 0.0;
      final double total = free + locked;

      if (total > 0.0001) {
        final price = prices[symbol] ?? 1.0;
        final value = total * price;
        totalValue += value;

        print(
          '    $symbol: ${total.toStringAsFixed(8)} × \$${price.toStringAsFixed(2)} = \$${value.toStringAsFixed(2)}',
        );
      }
    }

    return totalValue;
  } catch (e) {
    print('  ❌ Binance Spot Error: $e');
    return 0.0;
  }
}

Future<double> calculateBinanceFutures(BinanceApiClient client) async {
  try {
    final account = await client.getFuturesAccount();

    if (account?['totalWalletBalance'] != null) {
      final balance =
          double.tryParse(account?['totalWalletBalance']?.toString() ?? '0') ??
          0.0;
      print('  📋 Binance Futures: \$${balance.toStringAsFixed(2)} USDT');
      return balance;
    }

    return 0.0;
  } catch (e) {
    print('  ❌ Binance Futures Error: $e');
    return 0.0;
  }
}

Future<double> calculateGateSpot(
  GateApiClient client,
  Map<String, double> prices,
) async {
  try {
    final accounts = await client.getSpotAccounts();
    double totalValue = 0.0;
    int significantAssets = 0;

    print('  📋 Gate.io Spot Assets (significant holdings):');
    for (var account in accounts) {
      final symbol = account['currency'] as String;
      final double available =
          double.tryParse(account['available'].toString()) ?? 0.0;
      final double locked =
          double.tryParse(account['locked'].toString()) ?? 0.0;
      final double total = available + locked;

      if (total > 0.001) {
        // Only show significant amounts
        double price =
            prices[symbol] ?? 0.01; // Default very low price for unknown tokens

        // Special cases for known valuable tokens
        if (symbol == 'USDT' || symbol == 'USDC')
          price = 1.0;
        else if (symbol == 'BTC')
          price = prices['BTC'] ?? 45000.0;
        else if (symbol == 'ETH')
          price = prices['ETH'] ?? 2500.0;
        else if (symbol == 'BNB')
          price = prices['BNB'] ?? 300.0;
        else if (symbol == 'SOL')
          price = prices['SOL'] ?? 100.0;
        else if (symbol == 'GT')
          price = 5.0; // Gate token

        final value = total * price;
        totalValue += value;

        if (value > 0.10) {
          // Only show assets worth more than $0.10
          significantAssets++;
          print(
            '    $symbol: ${total.toStringAsFixed(6)} × \$${price.toStringAsFixed(4)} = \$${value.toStringAsFixed(2)}',
          );
        }
      }
    }

    if (significantAssets == 0) {
      print('    (Only dust amounts found)');
    }

    return totalValue;
  } catch (e) {
    print('  ❌ Gate.io Spot Error: $e');
    return 0.0;
  }
}

Future<double> calculateGateFutures(GateApiClient client) async {
  try {
    final account = await client.getFuturesAccount();

    if (account?['total'] != null) {
      final balance =
          double.tryParse(account?['total']?.toString() ?? '0') ?? 0.0;
      print('  📋 Gate.io Futures: \$${balance.toStringAsFixed(2)} USDT');
      return balance;
    }

    return 0.0;
  } catch (e) {
    print('  ❌ Gate.io Futures Error: $e');
    return 0.0;
  }
}
