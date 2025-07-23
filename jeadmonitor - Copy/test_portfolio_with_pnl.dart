import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import 'lib/services/binance_api_client.dart';
import 'lib/services/gate_api_client.dart';

void main() async {
  print('=== Portfolio Balance + PnL Calculator ===\n');

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

  // Get current prices
  final prices = await getCryptoPrices();
  print('📈 Current crypto prices loaded\n');

  double totalBalance = 0.0;
  double totalPnL = 0.0;

  print('🟡 === BINANCE ANALYSIS ===');
  final binanceResults = await analyzeBinance(binanceClient, prices);
  totalBalance += binanceResults['balance'] ?? 0.0;
  totalPnL += binanceResults['pnl'] ?? 0.0;

  print('\n🟢 === GATE.IO ANALYSIS ===');
  final gateResults = await analyzeGate(gateClient, prices);
  totalBalance += gateResults['balance'] ?? 0.0;
  totalPnL += gateResults['pnl'] ?? 0.0;

  print('\n' + '=' * 60);
  print('📊 PORTFOLIO SUMMARY');
  print('=' * 60);
  print('💰 Total Balance: \$${totalBalance.toStringAsFixed(2)}');
  print('📈 Total Unrealized PnL: \$${totalPnL.toStringAsFixed(2)}');
  print(
    '🎯 Total Portfolio Value: \$${(totalBalance + totalPnL).toStringAsFixed(2)}',
  );
  print('=' * 60);

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

Future<Map<String, double>> analyzeBinance(
  BinanceApiClient client,
  Map<String, double> prices,
) async {
  double totalBalance = 0.0;
  double totalPnL = 0.0;

  // Binance Spot
  print('📊 Binance Spot Account:');
  try {
    final account = await client.getAccountInfo();
    double spotValue = 0.0;

    for (var balance in account['balances']) {
      final symbol = balance['asset'] as String;
      final double free = double.tryParse(balance['free'].toString()) ?? 0.0;
      final double locked =
          double.tryParse(balance['locked'].toString()) ?? 0.0;
      final double total = free + locked;

      if (total > 0.0001) {
        final price = prices[symbol] ?? 1.0;
        final value = total * price;
        spotValue += value;

        print(
          '  ${symbol}: ${total.toStringAsFixed(8)} × \$${price.toStringAsFixed(2)} = \$${value.toStringAsFixed(2)}',
        );
      }
    }

    totalBalance += spotValue;
    print('  💰 Spot Total: \$${spotValue.toStringAsFixed(2)}');
  } catch (e) {
    print('  ❌ Spot Error: $e');
  }

  // Binance Futures
  print('\n📊 Binance Futures Account:');
  try {
    final account = await client.getFuturesAccount();

    if (account != null) {
      // Total Wallet Balance
      final walletBalance =
          double.tryParse(account['totalWalletBalance']?.toString() ?? '0') ??
          0.0;
      totalBalance += walletBalance;
      print('  💳 Total Wallet Balance: \$${walletBalance.toStringAsFixed(2)}');

      // Unrealized PnL
      final unrealizedPnL =
          double.tryParse(
            account['totalUnrealizedProfit']?.toString() ?? '0',
          ) ??
          0.0;
      totalPnL += unrealizedPnL;
      print('  📈 Total Unrealized PnL: \$${unrealizedPnL.toStringAsFixed(2)}');

      // Cross Wallet Balance
      final crossWalletBalance =
          double.tryParse(
            account['totalCrossWalletBalance']?.toString() ?? '0',
          ) ??
          0.0;
      print(
        '  🔄 Cross Wallet Balance: \$${crossWalletBalance.toStringAsFixed(2)}',
      );

      // Available Balance
      final availableBalance =
          double.tryParse(account['availableBalance']?.toString() ?? '0') ??
          0.0;
      print('  💵 Available Balance: \$${availableBalance.toStringAsFixed(2)}');

      // Individual asset details
      if (account['assets'] != null) {
        print('  📋 Asset Breakdown:');
        for (var asset in account['assets']) {
          final symbol = asset['asset'] as String;
          final walletBal =
              double.tryParse(asset['walletBalance']?.toString() ?? '0') ?? 0.0;
          final unrealizedProfit =
              double.tryParse(asset['unrealizedProfit']?.toString() ?? '0') ??
              0.0;

          if (walletBal > 0.001 || unrealizedProfit.abs() > 0.001) {
            print(
              '    ${symbol}: Wallet \$${walletBal.toStringAsFixed(2)}, PnL \$${unrealizedProfit.toStringAsFixed(2)}',
            );
          }
        }
      }

      // Positions with PnL
      if (account['positions'] != null) {
        print('  🎯 Active Positions:');
        for (var position in account['positions']) {
          final symbol = position['symbol'] as String;
          final positionAmt =
              double.tryParse(position['positionAmt']?.toString() ?? '0') ??
              0.0;
          final unrealizedProfit =
              double.tryParse(
                position['unrealizedProfit']?.toString() ?? '0',
              ) ??
              0.0;
          final percentage =
              double.tryParse(position['percentage']?.toString() ?? '0') ?? 0.0;

          if (positionAmt.abs() > 0.001 || unrealizedProfit.abs() > 0.001) {
            print(
              '    ${symbol}: Size ${positionAmt.toStringAsFixed(4)}, PnL \$${unrealizedProfit.toStringAsFixed(2)} (${percentage.toStringAsFixed(2)}%)',
            );
          }
        }
      }
    }
  } catch (e) {
    print('  ❌ Futures Error: $e');
  }

  print(
    '  🟡 Binance Summary: Balance \$${totalBalance.toStringAsFixed(2)}, PnL \$${totalPnL.toStringAsFixed(2)}',
  );
  return {'balance': totalBalance, 'pnl': totalPnL};
}

Future<Map<String, double>> analyzeGate(
  GateApiClient client,
  Map<String, double> prices,
) async {
  double totalBalance = 0.0;
  double totalPnL = 0.0;

  // Gate.io Spot
  print('📊 Gate.io Spot Account:');
  try {
    final accounts = await client.getSpotAccounts();
    double spotValue = 0.0;
    int significantAssets = 0;

    for (var account in accounts) {
      final symbol = account['currency'] as String;
      final double available =
          double.tryParse(account['available'].toString()) ?? 0.0;
      final double locked =
          double.tryParse(account['locked'].toString()) ?? 0.0;
      final double total = available + locked;

      if (total > 0.001) {
        double price = prices[symbol] ?? 0.01;

        // Known token prices
        if (symbol == 'USDT' || symbol == 'USDC')
          price = 1.0;
        else if (symbol == 'BTC')
          price = prices['BTC'] ?? 45000.0;
        else if (symbol == 'ETH')
          price = prices['ETH'] ?? 2500.0;
        else if (symbol == 'BNB')
          price = prices['BNB'] ?? 300.0;
        else if (symbol == 'GT')
          price = 5.0;

        final value = total * price;
        spotValue += value;

        if (value > 0.10) {
          significantAssets++;
          print(
            '  ${symbol}: ${total.toStringAsFixed(6)} × \$${price.toStringAsFixed(4)} = \$${value.toStringAsFixed(2)}',
          );
        }
      }
    }

    totalBalance += spotValue;
    print(
      '  💰 Spot Total: \$${spotValue.toStringAsFixed(2)} (${significantAssets} significant assets)',
    );
  } catch (e) {
    print('  ❌ Spot Error: $e');
  }

  // Gate.io Futures
  print('\n📊 Gate.io Futures Account:');
  try {
    final account = await client.getFuturesAccount();

    if (account != null) {
      // Total balance
      final total = double.tryParse(account['total']?.toString() ?? '0') ?? 0.0;
      totalBalance += total;
      print('  💳 Total Balance: \$${total.toStringAsFixed(2)}');

      // Available balance
      final available =
          double.tryParse(account['available']?.toString() ?? '0') ?? 0.0;
      print('  💵 Available: \$${available.toStringAsFixed(2)}');

      // Unrealized PnL
      final unrealizedPnL =
          double.tryParse(account['unrealised_pnl']?.toString() ?? '0') ?? 0.0;
      totalPnL += unrealizedPnL;
      print('  📈 Unrealized PnL: \$${unrealizedPnL.toStringAsFixed(2)}');

      // Position margin
      final positionMargin =
          double.tryParse(account['position_margin']?.toString() ?? '0') ?? 0.0;
      print('  🎯 Position Margin: \$${positionMargin.toStringAsFixed(2)}');

      // Order margin
      final orderMargin =
          double.tryParse(account['order_margin']?.toString() ?? '0') ?? 0.0;
      print('  📋 Order Margin: \$${orderMargin.toStringAsFixed(2)}');

      // Historical PnL and fees
      if (account['history'] != null) {
        final history = account['history'];
        final pnl = double.tryParse(history['pnl']?.toString() ?? '0') ?? 0.0;
        final fee = double.tryParse(history['fee']?.toString() ?? '0') ?? 0.0;
        final fund = double.tryParse(history['fund']?.toString() ?? '0') ?? 0.0;

        print('  📊 Historical Data:');
        print('    Realized PnL: \$${pnl.toStringAsFixed(2)}');
        print('    Total Fees: \$${fee.toStringAsFixed(2)}');
        print('    Funding: \$${fund.toStringAsFixed(2)}');
      }
    }
  } catch (e) {
    print('  ❌ Futures Error: $e');
  }

  print(
    '  🟢 Gate.io Summary: Balance \$${totalBalance.toStringAsFixed(2)}, PnL \$${totalPnL.toStringAsFixed(2)}',
  );
  return {'balance': totalBalance, 'pnl': totalPnL};
}
