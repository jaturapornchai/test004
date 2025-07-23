import 'dart:io';
import 'lib/services/binance_api_client.dart';
import 'lib/services/gate_api_client.dart';

void main() async {
  print('=== Crypto Portfolio Balance Test ===\n');

  // Initialize API clients
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

  print('🟡 Testing Binance APIs...\n');
  await testBinanceAPIs(binanceClient);

  print('\n🟢 Testing Gate.io APIs...\n');
  await testGateAPIs(gateClient);

  print('\n=== Test Complete ===');
  exit(0);
}

Future<void> testBinanceAPIs(BinanceApiClient client) async {
  print('📊 Testing Binance Spot Account...');
  try {
    final spotAccount = await client.getAccountInfo();
    print('✅ Spot Account Response: ${spotAccount.runtimeType}');

    if (spotAccount != null && spotAccount['balances'] != null) {
      final balances = spotAccount['balances'] as List;
      print('📈 Total balances count: ${balances.length}');

      double totalUSDValue = 0.0;
      int nonZeroBalances = 0;

      print('\n💰 Non-zero balances:');
      for (var balance in balances) {
        final double free = double.tryParse(balance['free'].toString()) ?? 0.0;
        final double locked =
            double.tryParse(balance['locked'].toString()) ?? 0.0;
        final double total = free + locked;

        if (total > 0.0001) {
          // Filter out dust
          nonZeroBalances++;
          print(
            '  ${balance['asset']}: ${total.toStringAsFixed(8)} (Free: ${free.toStringAsFixed(8)}, Locked: ${locked.toStringAsFixed(8)})',
          );

          // Mock USD conversion (you would need price API for real values)
          if (balance['asset'] == 'USDT' ||
              balance['asset'] == 'BUSD' ||
              balance['asset'] == 'USDC') {
            totalUSDValue += total;
          } else if (balance['asset'] == 'BTC') {
            totalUSDValue += total * 45000; // Mock BTC price
          } else if (balance['asset'] == 'ETH') {
            totalUSDValue += total * 2500; // Mock ETH price
          } else if (balance['asset'] == 'BNB') {
            totalUSDValue += total * 300; // Mock BNB price
          } else {
            totalUSDValue += total * 1; // Mock other coins as $1
          }
        }
      }

      print('📊 Summary: ${nonZeroBalances} assets with balance');
      print(
        '💵 Estimated total USD value: \$${totalUSDValue.toStringAsFixed(2)}',
      );
    }
  } catch (e) {
    print('❌ Binance Spot Error: $e');
  }

  print('\n📊 Testing Binance Futures Account...');
  try {
    final futuresAccount = await client.getFuturesAccount();
    print('✅ Futures Account Response: ${futuresAccount.runtimeType}');

    if (futuresAccount != null) {
      // Check for different possible response structures
      if (futuresAccount['assets'] != null) {
        final assets = futuresAccount['assets'] as List;
        print('📈 Total futures assets: ${assets.length}');

        double totalUSDValue = 0.0;
        int nonZeroAssets = 0;

        print('\n💰 Non-zero futures balances:');
        for (var asset in assets) {
          final double walletBalance =
              double.tryParse(asset['walletBalance'].toString()) ?? 0.0;
          final double unrealizedProfit =
              double.tryParse(asset['unrealizedProfit'].toString()) ?? 0.0;

          if (walletBalance > 0.0001 || unrealizedProfit.abs() > 0.0001) {
            nonZeroAssets++;
            print(
              '  ${asset['asset']}: Wallet: ${walletBalance.toStringAsFixed(8)}, PnL: ${unrealizedProfit.toStringAsFixed(8)}',
            );

            if (asset['asset'] == 'USDT') {
              totalUSDValue += walletBalance;
            }
          }
        }

        print('📊 Summary: ${nonZeroAssets} futures assets');
        print(
          '💵 Estimated futures USD value: \$${totalUSDValue.toStringAsFixed(2)}',
        );
      }

      if (futuresAccount['totalWalletBalance'] != null) {
        final totalBalance =
            double.tryParse(futuresAccount['totalWalletBalance'].toString()) ??
            0.0;
        print('💰 Total Wallet Balance: \$${totalBalance.toStringAsFixed(2)}');
      }
    }
  } catch (e) {
    print('❌ Binance Futures Error: $e');
  }
}

Future<void> testGateAPIs(GateApiClient client) async {
  print('📊 Testing Gate.io Spot Account...');
  try {
    final spotAccounts = await client.getSpotAccounts();
    print('✅ Spot Accounts Response: ${spotAccounts.runtimeType}');

    if (spotAccounts != null && spotAccounts is List) {
      print('📈 Total spot accounts: ${spotAccounts.length}');

      double totalUSDValue = 0.0;
      int nonZeroAccounts = 0;

      print('\n💰 Non-zero spot balances:');
      for (var account in spotAccounts) {
        final double available =
            double.tryParse(account['available'].toString()) ?? 0.0;
        final double locked =
            double.tryParse(account['locked'].toString()) ?? 0.0;
        final double total = available + locked;

        if (total > 0.0001) {
          // Filter out dust
          nonZeroAccounts++;
          print(
            '  ${account['currency']}: ${total.toStringAsFixed(8)} (Available: ${available.toStringAsFixed(8)}, Locked: ${locked.toStringAsFixed(8)})',
          );

          // Mock USD conversion
          if (account['currency'] == 'USDT' || account['currency'] == 'USDC') {
            totalUSDValue += total;
          } else if (account['currency'] == 'BTC') {
            totalUSDValue += total * 45000; // Mock BTC price
          } else if (account['currency'] == 'ETH') {
            totalUSDValue += total * 2500; // Mock ETH price
          } else if (account['currency'] == 'GT') {
            totalUSDValue += total * 5; // Mock GT price
          } else {
            totalUSDValue += total * 1; // Mock other coins as $1
          }
        }
      }

      print('📊 Summary: ${nonZeroAccounts} spot accounts with balance');
      print(
        '💵 Estimated total USD value: \$${totalUSDValue.toStringAsFixed(2)}',
      );
    }
  } catch (e) {
    print('❌ Gate.io Spot Error: $e');
  }

  print('\n📊 Testing Gate.io Futures Account...');
  try {
    final futuresAccount = await client.getFuturesAccount();
    print('✅ Futures Account Response: ${futuresAccount.runtimeType}');
    print('📄 Raw response: $futuresAccount');

    if (futuresAccount != null) {
      // Check different possible fields
      final fields = [
        'total',
        'available',
        'balance',
        'equity',
        'unrealized_pnl',
      ];
      for (String field in fields) {
        if (futuresAccount[field] != null) {
          final value =
              double.tryParse(futuresAccount[field].toString()) ?? 0.0;
          print('💰 $field: \$${value.toStringAsFixed(2)}');
        }
      }
    }
  } catch (e) {
    print('❌ Gate.io Futures Error: $e');
  }
}
