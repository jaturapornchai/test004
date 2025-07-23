import 'dart:io';
import 'test_gate_new_client.dart';

void main() async {
  print('=== Testing New Gate.io API Client ===\n');

  final client = GateApiClientNew(
    apiKey: '472667d2eebf07b206f8599bf35bbcb1',
    apiSecret:
        'e24966020b7241e521966b87deb401d6dac0f90a1ccbf3f85db54ecc09faf33f',
  );

  print('🟢 Testing Gate.io Spot Account...\n');
  try {
    final spotAccounts = await client.getSpotAccounts();
    print('✅ Gate.io Spot Success!');
    print('📊 Found ${spotAccounts.length} spot accounts');

    double totalValue = 0.0;
    int nonZeroAccounts = 0;

    print('\n💰 Significant balances:');
    for (var account in spotAccounts) {
      final symbol = account['currency'] as String;
      final double available =
          double.tryParse(account['available'].toString()) ?? 0.0;
      final double locked =
          double.tryParse(account['locked'].toString()) ?? 0.0;
      final double total = available + locked;

      if (total > 0.001) {
        nonZeroAccounts++;
        double estimatedValue = total;
        if (symbol == 'USDT' || symbol == 'USDC') {
          estimatedValue = total * 1.0;
          totalValue += estimatedValue;
        } else if (symbol == 'BTC') {
          estimatedValue = total * 45000;
          totalValue += estimatedValue;
        } else if (symbol == 'GT') {
          estimatedValue = total * 5.0;
          totalValue += estimatedValue;
        }

        if (estimatedValue > 0.10) {
          print(
            '  ${symbol}: ${total.toStringAsFixed(6)} ≈ \$${estimatedValue.toStringAsFixed(2)}',
          );
        }
      }
    }

    print('\n📊 Summary: ${nonZeroAccounts} accounts with balance');
    print('💵 Estimated total value: \$${totalValue.toStringAsFixed(2)}');
  } catch (e) {
    print('❌ Gate.io Spot Test Failed: $e');
  }

  print('\n🟢 Testing Gate.io Futures Account...\n');
  try {
    final futuresAccount = await client.getFuturesAccount();

    if (futuresAccount != null) {
      print('✅ Gate.io Futures Success!');

      final total =
          double.tryParse(futuresAccount['total']?.toString() ?? '0') ?? 0.0;
      final available =
          double.tryParse(futuresAccount['available']?.toString() ?? '0') ??
          0.0;
      final unrealizedPnL =
          double.tryParse(
            futuresAccount['unrealised_pnl']?.toString() ?? '0',
          ) ??
          0.0;

      print('💰 Total Balance: \$${total.toStringAsFixed(2)}');
      print('💵 Available: \$${available.toStringAsFixed(2)}');
      print('📈 Unrealized PnL: \$${unrealizedPnL.toStringAsFixed(2)}');

      if (futuresAccount['history'] != null) {
        final history = futuresAccount['history'];
        final realizedPnL =
            double.tryParse(history['pnl']?.toString() ?? '0') ?? 0.0;
        final totalFees =
            double.tryParse(history['fee']?.toString() ?? '0') ?? 0.0;

        print('📊 Historical Data:');
        print('  Realized PnL: \$${realizedPnL.toStringAsFixed(2)}');
        print('  Total Fees: \$${totalFees.toStringAsFixed(2)}');
      }
    } else {
      print('❌ Gate.io Futures returned null');
    }
  } catch (e) {
    print('❌ Gate.io Futures Test Failed: $e');
  }

  print('\n=== Test Complete ===');
  exit(0);
}
