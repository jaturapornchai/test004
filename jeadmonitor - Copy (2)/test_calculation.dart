import 'dart:convert';
import 'dart:io';
import 'package:crypto/crypto.dart';

void main() async {
  // Test the calculation precision
  double binanceSpot = 0.00;
  double binanceFutures = 197.82239510;
  double binancePnL = -0.13261379;

  double gateSpot = 0.00;
  double gateFutures = 1206.063721773725;
  double gatePnL = 28.059010038774;

  print('=== Current Values ===');
  print('Binance Spot: \$${binanceSpot.toStringAsFixed(2)}');
  print('Binance Futures: \$${binanceFutures.toStringAsFixed(2)}');
  print('Binance PnL: \$${binancePnL.toStringAsFixed(2)}');
  print(
    'Binance Total: \$${(binanceSpot + binanceFutures).toStringAsFixed(2)}',
  );

  print('\nGate.io Spot: \$${gateSpot.toStringAsFixed(2)}');
  print('Gate.io Futures: \$${gateFutures.toStringAsFixed(2)}');
  print('Gate.io PnL: \$${gatePnL.toStringAsFixed(2)}');
  print('Gate.io Total: \$${(gateSpot + gateFutures).toStringAsFixed(2)}');

  double portfolioTotal = binanceSpot + binanceFutures + gateSpot + gateFutures;
  print('\n=== Portfolio Total ===');
  print('Portfolio Total: \$${portfolioTotal.toStringAsFixed(2)}');
  print('Expected: \$${(197.82 + 1206.06).toStringAsFixed(2)}');

  // Test rounding difference
  print('\nBinance rounded: \$${binanceFutures.toStringAsFixed(2)}');
  print('Gate.io rounded: \$${gateFutures.toStringAsFixed(2)}');
  print(
    'Sum of rounded: \$${(double.parse(binanceFutures.toStringAsFixed(2)) + double.parse(gateFutures.toStringAsFixed(2))).toStringAsFixed(2)}',
  );
}
