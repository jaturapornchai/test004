import '../services/binance_api_client.dart';
import '../services/gate_api_client.dart';

class RealizedPnLService {
  final BinanceApiClient binanceClient;
  final GateApiClient gateClient;

  RealizedPnLService({required this.binanceClient, required this.gateClient});

  // Get today's realized P&L from closed positions
  Future<Map<String, dynamic>> getTodayRealizedPnL() async {
    try {
      final results = await Future.wait([
        _getBinanceSpotRealizedPnL(),
        _getBinanceFuturesRealizedPnL(),
        _getGateSpotRealizedPnL(),
        _getGateFuturesRealizedPnL(),
      ]);

      final binanceSpotPnL = results[0];
      final binanceFuturesPnL = results[1];
      final gateSpotPnL = results[2];
      final gateFuturesPnL = results[3];

      // Calculate total realized P&L
      final totalRealizedPnL =
          binanceSpotPnL['realizedPnL']! +
          binanceFuturesPnL['realizedPnL']! +
          gateSpotPnL['realizedPnL']! +
          gateFuturesPnL['realizedPnL']!;

      return {
        'totalRealizedPnL': totalRealizedPnL,
        'binanceSpot': binanceSpotPnL,
        'binanceFutures': binanceFuturesPnL,
        'gateSpot': gateSpotPnL,
        'gateFutures': gateFuturesPnL,
      };
    } catch (e) {
      print('Error calculating today realized P&L: $e');
      return {
        'totalRealizedPnL': 0.0,
        'binanceSpot': {'realizedPnL': 0.0, 'tradeCount': 0},
        'binanceFutures': {'realizedPnL': 0.0, 'tradeCount': 0},
        'gateSpot': {'realizedPnL': 0.0, 'tradeCount': 0},
        'gateFutures': {'realizedPnL': 0.0, 'tradeCount': 0},
      };
    }
  }

  // Calculate Binance Spot realized P&L from today's trades
  Future<Map<String, double>> _getBinanceSpotRealizedPnL() async {
    try {
      final trades = await binanceClient.getTodaySpotTradeHistory();

      double totalRealizedPnL = 0.0;
      Map<String, List<Map<String, dynamic>>> symbolTrades = {};

      // Group trades by symbol
      for (var trade in trades) {
        final symbol = trade['symbol'] as String;
        if (!symbolTrades.containsKey(symbol)) {
          symbolTrades[symbol] = [];
        }
        symbolTrades[symbol]!.add(trade);
      }

      // Calculate P&L for each symbol
      for (var symbol in symbolTrades.keys) {
        final symbolTradeList = symbolTrades[symbol]!;

        // Calculate weighted average buy/sell prices and quantities
        double totalBuyQty = 0.0;
        double totalBuyValue = 0.0;
        double totalSellQty = 0.0;
        double totalSellValue = 0.0;

        for (var trade in symbolTradeList) {
          final qty = double.tryParse(trade['qty']?.toString() ?? '0') ?? 0.0;
          final price =
              double.tryParse(trade['price']?.toString() ?? '0') ?? 0.0;
          final isBuyer = trade['isBuyer'] as bool;

          if (isBuyer) {
            totalBuyQty += qty;
            totalBuyValue += qty * price;
          } else {
            totalSellQty += qty;
            totalSellValue += qty * price;
          }
        }

        // Calculate P&L for completed trades (min of buy/sell quantities)
        final completedQty = totalBuyQty < totalSellQty
            ? totalBuyQty
            : totalSellQty;
        if (completedQty > 0) {
          final avgBuyPrice = totalBuyValue / totalBuyQty;
          final avgSellPrice = totalSellValue / totalSellQty;
          final symbolPnL = completedQty * (avgSellPrice - avgBuyPrice);
          totalRealizedPnL += symbolPnL;
        }
      }

      return {
        'realizedPnL': totalRealizedPnL,
        'tradeCount': trades.length.toDouble(),
      };
    } catch (e) {
      print('Error calculating Binance Spot realized P&L: $e');
      return {'realizedPnL': 0.0, 'tradeCount': 0.0};
    }
  }

  // Calculate Binance Futures realized P&L from today's trades
  Future<Map<String, double>> _getBinanceFuturesRealizedPnL() async {
    try {
      final trades = await binanceClient.getTodayFuturesTradeHistory();

      double totalRealizedPnL = 0.0;

      for (var trade in trades) {
        // In futures, realizedPnl field directly shows the P&L for this trade
        final realizedPnl =
            double.tryParse(trade['realizedPnl']?.toString() ?? '0') ?? 0.0;
        totalRealizedPnL += realizedPnl;
      }

      return {
        'realizedPnL': totalRealizedPnL,
        'tradeCount': trades.length.toDouble(),
      };
    } catch (e) {
      print('Error calculating Binance Futures realized P&L: $e');
      return {'realizedPnL': 0.0, 'tradeCount': 0.0};
    }
  }

  // Calculate Gate.io Spot realized P&L from today's trades
  Future<Map<String, double>> _getGateSpotRealizedPnL() async {
    try {
      final trades = await gateClient.getTodaySpotTradeHistory();

      double totalRealizedPnL = 0.0;
      Map<String, List<Map<String, dynamic>>> symbolTrades = {};

      // Group trades by currency pair
      for (var trade in trades) {
        final symbol = trade['currency_pair'] as String;
        if (!symbolTrades.containsKey(symbol)) {
          symbolTrades[symbol] = [];
        }
        symbolTrades[symbol]!.add(trade);
      }

      // Calculate P&L for each symbol
      for (var symbol in symbolTrades.keys) {
        final symbolTradeList = symbolTrades[symbol]!;

        double totalBuyQty = 0.0;
        double totalBuyValue = 0.0;
        double totalSellQty = 0.0;
        double totalSellValue = 0.0;

        for (var trade in symbolTradeList) {
          final amount =
              double.tryParse(trade['amount']?.toString() ?? '0') ?? 0.0;
          final price =
              double.tryParse(trade['price']?.toString() ?? '0') ?? 0.0;
          final side = trade['side'] as String;

          if (side == 'buy') {
            totalBuyQty += amount;
            totalBuyValue += amount * price;
          } else {
            totalSellQty += amount;
            totalSellValue += amount * price;
          }
        }

        // Calculate P&L for completed trades
        final completedQty = totalBuyQty < totalSellQty
            ? totalBuyQty
            : totalSellQty;
        if (completedQty > 0) {
          final avgBuyPrice = totalBuyValue / totalBuyQty;
          final avgSellPrice = totalSellValue / totalSellQty;
          final symbolPnL = completedQty * (avgSellPrice - avgBuyPrice);
          totalRealizedPnL += symbolPnL;
        }
      }

      return {
        'realizedPnL': totalRealizedPnL,
        'tradeCount': trades.length.toDouble(),
      };
    } catch (e) {
      print('Error calculating Gate.io Spot realized P&L: $e');
      return {'realizedPnL': 0.0, 'tradeCount': 0.0};
    }
  }

  // Calculate Gate.io Futures realized P&L from today's trades
  Future<Map<String, double>> _getGateFuturesRealizedPnL() async {
    try {
      final trades = await gateClient.getTodayFuturesTradeHistory();

      double totalRealizedPnL = 0.0;

      for (var trade in trades) {
        // Calculate P&L based on point_fee or profit field if available
        final pointFee =
            double.tryParse(trade['point_fee']?.toString() ?? '0') ?? 0.0;
        final profit =
            double.tryParse(trade['profit']?.toString() ?? '0') ?? 0.0;

        // Use profit field if available, otherwise calculate from point_fee
        final tradePnL = profit != 0.0 ? profit : pointFee;
        totalRealizedPnL += tradePnL;
      }

      return {
        'realizedPnL': totalRealizedPnL,
        'tradeCount': trades.length.toDouble(),
      };
    } catch (e) {
      print('Error calculating Gate.io Futures realized P&L: $e');
      return {'realizedPnL': 0.0, 'tradeCount': 0.0};
    }
  }
}
