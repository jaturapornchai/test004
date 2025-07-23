import 'dart:convert';
import 'dart:io';
import 'package:crypto/crypto.dart';

void main() async {
  print('=== Testing Current Values vs Correct Values ===');

  // Current values from our test
  print('\n--- Our Test Results ---');
  print('Gate.io Total: 1206.063721773725 USDT');
  print('Gate.io PnL: 28.059010038774 USDT');
  print('Binance Balance: 197.82239510 USDT');
  print('Binance PnL: -0.13261379 USDT');

  // Correct values from user's screenshot
  print('\n--- Correct Values (from screenshot) ---');
  print('Gate.io Account Assets: 1235.31 USDT');
  print('Gate.io Wallet Assets: 1206.06 USDT');
  print('Gate.io Unrealized PnL: +29.24 USDT');
  print('Binance Balance: 197.8223 USDT');
  print('Binance Unrealized PnL: 1.2707 USDT');

  // Analysis
  print('\n--- Analysis ---');
  print(
    'Gate.io: Account Assets (1235.31) = Wallet Assets (1206.06) + PnL (29.24)',
  );
  print('Calculation: 1206.06 + 29.24 = 1235.30 ≈ 1235.31 ✓');

  print('Binance: Total should be Balance + PnL');
  print('Calculation: 197.8223 + 1.2707 = 199.093 USDT');

  print('\n--- Issues Found ---');
  print('1. Gate.io PnL should be 29.24, not 28.06');
  print('2. Binance PnL should be 1.2707, not -0.13');
  print(
    '3. Need to use Account Assets (1235.31) for Gate.io total, not Wallet Assets',
  );
}
