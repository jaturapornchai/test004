import '../services/binance_api_client.dart';
import '../services/gate_api_client.dart';

class DailyPnLService {
  final BinanceApiClient binanceClient;
  final GateApiClient gateClient;

  DailyPnLService({required this.binanceClient, required this.gateClient});

  // Get portfolio daily P&L from API data
  Future<Map<String, dynamic>> getDailyPortfolioPnL() async {
    try {
      final results = await Future.wait([
        _getBinanceSpotPnL(),
        _getBinanceFuturesPnL(),
        _getGateSpotPnL(),
        _getGateFuturesPnL(),
      ]);

      final binanceSpotPnL = results[0];
      final binanceFuturesPnL = results[1];
      final gateSpotPnL = results[2];
      final gateFuturesPnL = results[3];

      // Calculate total daily P&L
      final totalDailyPnL =
          binanceSpotPnL['totalPnL']! +
          binanceFuturesPnL['totalPnL']! +
          gateSpotPnL['totalPnL']! +
          gateFuturesPnL['totalPnL']!;

      final totalPortfolioValue =
          binanceSpotPnL['currentValue']! +
          binanceFuturesPnL['currentValue']! +
          gateSpotPnL['currentValue']! +
          gateFuturesPnL['currentValue']!;

      final totalDailyPnLPercent = totalPortfolioValue > 0
          ? (totalDailyPnL / (totalPortfolioValue - totalDailyPnL)) * 100
          : 0.0;

      return {
        'totalDailyPnL': totalDailyPnL,
        'totalDailyPnLPercent': totalDailyPnLPercent,
        'binanceSpot': binanceSpotPnL,
        'binanceFutures': binanceFuturesPnL,
        'gateSpot': gateSpotPnL,
        'gateFutures': gateFuturesPnL,
        'totalPortfolioValue': totalPortfolioValue,
      };
    } catch (e) {
      print('Error calculating daily P&L: $e');
      return {
        'totalDailyPnL': 0.0,
        'totalDailyPnLPercent': 0.0,
        'binanceSpot': {
          'totalPnL': 0.0,
          'currentValue': 0.0,
          'changePercent': 0.0,
        },
        'binanceFutures': {
          'totalPnL': 0.0,
          'currentValue': 0.0,
          'changePercent': 0.0,
        },
        'gateSpot': {
          'totalPnL': 0.0,
          'currentValue': 0.0,
          'changePercent': 0.0,
        },
        'gateFutures': {
          'totalPnL': 0.0,
          'currentValue': 0.0,
          'changePercent': 0.0,
        },
        'totalPortfolioValue': 0.0,
      };
    }
  }

  // Calculate Binance Spot 24hr P&L
  Future<Map<String, double>> _getBinanceSpotPnL() async {
    try {
      final accountInfo = await binanceClient.getAccountInfo();
      final ticker24hr = await binanceClient.get24hrTicker();

      // Create price change map
      Map<String, double> priceChangePercent = {};
      for (var ticker in ticker24hr) {
        final symbol = ticker['symbol'] as String;
        final changePercent =
            double.tryParse(ticker['priceChangePercent']?.toString() ?? '0') ??
            0.0;
        if (symbol.endsWith('USDT')) {
          final asset = symbol.replaceAll('USDT', '');
          priceChangePercent[asset] = changePercent;
        }
      }

      double totalCurrentValue = 0.0;
      double totalPnL = 0.0;

      for (var balance in accountInfo['balances']) {
        final double free = double.tryParse(balance['free']) ?? 0.0;
        final double locked = double.tryParse(balance['locked']) ?? 0.0;
        final double total = free + locked;
        final String asset = balance['asset'];

        if (total > 0) {
          double currentUsdValue = 0.0;

          if (asset == 'USDT' || asset == 'USDC' || asset == 'BUSD') {
            // Stablecoins - no price change
            currentUsdValue = total;
          } else if (priceChangePercent.containsKey(asset)) {
            // Get current price from ticker
            final currentTicker = ticker24hr.firstWhere(
              (t) => t['symbol'] == '${asset}USDT',
              orElse: () => {'lastPrice': '0'},
            );
            final currentPrice =
                double.tryParse(
                  currentTicker['lastPrice']?.toString() ?? '0',
                ) ??
                0.0;
            currentUsdValue = total * currentPrice;

            // Calculate yesterday's value using price change
            final changePercent = priceChangePercent[asset]!;
            final yesterdayPrice = currentPrice / (1 + changePercent / 100);
            final yesterdayValue = total * yesterdayPrice;

            totalPnL += (currentUsdValue - yesterdayValue);
          }

          if (currentUsdValue > 0.01) {
            totalCurrentValue += currentUsdValue;
          }
        }
      }

      final weightedChangePercent = totalCurrentValue > 0
          ? (totalPnL / (totalCurrentValue - totalPnL)) * 100
          : 0.0;

      return {
        'totalPnL': totalPnL,
        'currentValue': totalCurrentValue,
        'changePercent': weightedChangePercent,
      };
    } catch (e) {
      print('Error calculating Binance Spot P&L: $e');
      return {'totalPnL': 0.0, 'currentValue': 0.0, 'changePercent': 0.0};
    }
  }

  // Calculate Binance Futures 24hr P&L
  Future<Map<String, double>> _getBinanceFuturesPnL() async {
    try {
      final accountInfo = await binanceClient.getFuturesAccount();

      final totalWalletBalance =
          double.tryParse(
            accountInfo?['totalWalletBalance']?.toString() ?? '0',
          ) ??
          0.0;

      final totalUnrealizedProfit =
          double.tryParse(
            accountInfo?['totalUnrealizedProfit']?.toString() ?? '0',
          ) ??
          0.0;

      // For futures, we can use the unrealized PnL as daily change
      // since it represents current open positions profit/loss
      final currentValue = totalWalletBalance;

      return {
        'totalPnL': totalUnrealizedProfit,
        'currentValue': currentValue,
        'changePercent': currentValue > 0
            ? (totalUnrealizedProfit / currentValue) * 100
            : 0.0,
      };
    } catch (e) {
      print('Error calculating Binance Futures P&L: $e');
      return {'totalPnL': 0.0, 'currentValue': 0.0, 'changePercent': 0.0};
    }
  }

  // Calculate Gate.io Spot 24hr P&L
  Future<Map<String, double>> _getGateSpotPnL() async {
    try {
      final accounts = await gateClient.getSpotAccounts();
      final tickers = await gateClient.getAllTickers();

      // Create price change map from tickers
      Map<String, Map<String, double>> tickerData = {};
      for (var ticker in tickers) {
        final symbol = ticker['currency_pair'] as String;
        final changePercent =
            double.tryParse(ticker['change_percentage']?.toString() ?? '0') ??
            0.0;
        final currentPrice =
            double.tryParse(ticker['last']?.toString() ?? '0') ?? 0.0;

        if (symbol.endsWith('_USDT')) {
          final asset = symbol.replaceAll('_USDT', '');
          tickerData[asset] = {
            'currentPrice': currentPrice,
            'changePercent': changePercent,
          };
        }
      }

      double totalCurrentValue = 0.0;
      double totalPnL = 0.0;

      for (var account in accounts) {
        final double available = double.tryParse(account['available']) ?? 0.0;
        final double locked = double.tryParse(account['locked']) ?? 0.0;
        final double total = available + locked;
        final String currency = account['currency'] ?? '';

        if (total > 0) {
          double currentUsdValue = 0.0;

          if (currency == 'USDT' || currency == 'USDC') {
            // Stablecoins - no price change
            currentUsdValue = total;
          } else if (tickerData.containsKey(currency)) {
            final data = tickerData[currency]!;
            final currentPrice = data['currentPrice']!;
            final changePercent = data['changePercent']!;

            currentUsdValue = total * currentPrice;

            // Calculate yesterday's value
            final yesterdayPrice = currentPrice / (1 + changePercent / 100);
            final yesterdayValue = total * yesterdayPrice;

            totalPnL += (currentUsdValue - yesterdayValue);
          }

          if (currentUsdValue > 0.01) {
            totalCurrentValue += currentUsdValue;
          }
        }
      }

      final weightedChangePercent = totalCurrentValue > 0
          ? (totalPnL / (totalCurrentValue - totalPnL)) * 100
          : 0.0;

      return {
        'totalPnL': totalPnL,
        'currentValue': totalCurrentValue,
        'changePercent': weightedChangePercent,
      };
    } catch (e) {
      print('Error calculating Gate.io Spot P&L: $e');
      return {'totalPnL': 0.0, 'currentValue': 0.0, 'changePercent': 0.0};
    }
  }

  // Calculate Gate.io Futures 24hr P&L
  Future<Map<String, double>> _getGateFuturesPnL() async {
    try {
      final accountInfo = await gateClient.getFuturesAccount();

      if (accountInfo != null) {
        final totalValue =
            double.tryParse(accountInfo['total']?.toString() ?? '0') ?? 0.0;
        final unrealizedPnL =
            double.tryParse(accountInfo['unrealised_pnl']?.toString() ?? '0') ??
            0.0;

        // Similar to Binance Futures, use unrealized PnL as daily change
        return {
          'totalPnL': unrealizedPnL,
          'currentValue': totalValue,
          'changePercent': totalValue > 0
              ? (unrealizedPnL / totalValue) * 100
              : 0.0,
        };
      }

      return {'totalPnL': 0.0, 'currentValue': 0.0, 'changePercent': 0.0};
    } catch (e) {
      print('Error calculating Gate.io Futures P&L: $e');
      return {'totalPnL': 0.0, 'currentValue': 0.0, 'changePercent': 0.0};
    }
  }
}
