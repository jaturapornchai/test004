import 'dart:async';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../services/binance_api_client.dart';
import '../services/gate_api_client.dart';
import '../services/exchange_rate_service.dart';
import '../config/api_config.dart';

part 'portfolio_state.dart';

class PortfolioCubit extends Cubit<PortfolioState> {
  final BinanceApiClient _binanceClient;
  final GateApiClient _gateClient;
  final ExchangeRateService _exchangeRateService;
  Timer? _refreshTimer;

  PortfolioCubit()
    : _binanceClient = BinanceApiClient(
        apiKey: ApiConfig.binanceApiKey,
        apiSecret: ApiConfig.binanceApiSecret,
      ),
      _gateClient = GateApiClient(
        apiKey: ApiConfig.gateApiKey,
        apiSecret: ApiConfig.gateApiSecret,
      ),
      _exchangeRateService = ExchangeRateService(),
      super(const PortfolioState());

  Future<void> initialize() async {
    emit(state.copyWith(status: PortfolioStatus.loading));

    try {
      // Load exchange rate once at startup
      final exchangeRate = await _exchangeRateService.getUsdToThbRate();

      // Load all portfolio data
      await _loadAllData(exchangeRate);

      // Start refresh timer only after initial load is complete
      _startRefreshTimer();
    } catch (e) {
      emit(state.copyWith(status: PortfolioStatus.error, error: e.toString()));
    }
  }

  Future<void> _loadAllData(double exchangeRate) async {
    try {
      // Load all data in parallel
      final results = await Future.wait([
        _loadBinanceSpotAssets(),
        _loadBinanceFuturesAssets(),
        _loadGateSpotAssets(),
        _loadGateFuturesAssets(),
        _loadBtcPrice(),
      ]);

      final binanceSpot = results[0] as Map<String, double>;
      final binanceFutures = results[1] as Map<String, double>;
      final gateSpot = results[2] as double;
      final gateFutures = results[3] as Map<String, double>;
      final btcPrice = results[4] as double;

      emit(
        state.copyWith(
          status: PortfolioStatus.loaded,
          binanceSpotValue: binanceSpot['value']!,
          binanceFuturesValue: binanceFutures['value']!,
          binanceFuturesPnL: binanceFutures['pnl']!,
          gateSpotValue: gateSpot,
          gateFuturesValue: gateFutures['value']!,
          gateFuturesPnL: gateFutures['pnl']!,
          usdToThbRate: exchangeRate,
          btcPrice: btcPrice,
        ),
      );
    } catch (e) {
      emit(state.copyWith(status: PortfolioStatus.error, error: e.toString()));
    }
  }

  void _startRefreshTimer() {
    _refreshTimer = Timer.periodic(const Duration(seconds: 5), (timer) {
      _refreshData();
    });
  }

  Future<void> _refreshData() async {
    if (state.status == PortfolioStatus.loading) return;

    try {
      // Use existing exchange rate for refresh
      await _loadAllData(state.usdToThbRate);
    } catch (e) {
      emit(state.copyWith(status: PortfolioStatus.error, error: e.toString()));
    }
  }

  Future<Map<String, double>> _loadBinanceSpotAssets() async {
    try {
      final accountInfo = await _binanceClient.getAccountInfo();
      final priceData = await _binanceClient.getAllPrices();

      // Build price map
      final prices = <String, double>{};
      for (var price in priceData) {
        final symbol = price['symbol'] as String;
        final priceValue = double.tryParse(price['price']) ?? 0.0;
        if (symbol.endsWith('USDT')) {
          final asset = symbol.replaceAll('USDT', '');
          prices[asset] = priceValue;
        }
      }

      // Calculate total spot value
      double totalSpotValue = 0.0;
      for (var balance in accountInfo['balances']) {
        final double free = double.tryParse(balance['free']) ?? 0.0;
        final double locked = double.tryParse(balance['locked']) ?? 0.0;
        final double total = free + locked;
        final String asset = balance['asset'];

        if (total > 0) {
          double usdValue = 0.0;
          if (asset == 'USDT' || asset == 'USDC' || asset == 'BUSD') {
            usdValue = total;
          } else if (prices.containsKey(asset)) {
            usdValue = total * prices[asset]!;
          }

          if (usdValue > 0.01) {
            totalSpotValue += usdValue;
          }
        }
      }

      return {'value': totalSpotValue};
    } catch (e) {
      print('Error loading Binance Spot assets: $e');
      return {'value': 0.0};
    }
  }

  Future<Map<String, double>> _loadBinanceFuturesAssets() async {
    try {
      final accountInfo = await _binanceClient.getFuturesAccount();

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
      final balanceWithoutPnL = totalWalletBalance + totalUnrealizedProfit;
      return {'value': balanceWithoutPnL, 'pnl': totalUnrealizedProfit};
    } catch (e) {
      print('Error loading Binance Futures assets: $e');
      return {'value': 0.0, 'pnl': 0.0};
    }
  }

  Future<double> _loadGateSpotAssets() async {
    try {
      final accounts = await _gateClient.getSpotAccounts();
      final tickers = await _gateClient.getAllTickers();

      // Build price map
      final prices = <String, double>{};
      for (var ticker in tickers) {
        final symbol = ticker['currency_pair'] as String;
        final lastPrice =
            double.tryParse(ticker['last']?.toString() ?? '0') ?? 0.0;
        if (symbol.endsWith('_USDT')) {
          final asset = symbol.replaceAll('_USDT', '');
          prices[asset] = lastPrice;
        }
      }

      // Calculate total spot value
      double totalSpotValue = 0.0;
      for (var account in accounts) {
        final double available = double.tryParse(account['available']) ?? 0.0;
        final double locked = double.tryParse(account['locked']) ?? 0.0;
        final double total = available + locked;
        final String currency = account['currency'] ?? '';

        if (total > 0) {
          double usdValue = 0.0;
          if (currency == 'USDT' || currency == 'USDC') {
            usdValue = total;
          } else if (prices.containsKey(currency)) {
            usdValue = total * prices[currency]!;
          }

          if (usdValue > 0.01) {
            totalSpotValue += usdValue;
          }
        }
      }

      return totalSpotValue;
    } catch (e) {
      print('Error loading Gate.io Spot assets: $e');
      return 0.0;
    }
  }

  Future<Map<String, double>> _loadGateFuturesAssets() async {
    try {
      final accountInfo = await _gateClient.getFuturesAccount();

      if (accountInfo != null) {
        final total =
            double.tryParse(accountInfo['total']?.toString() ?? '0') ?? 0.0;
        final unrealizedPnL =
            double.tryParse(accountInfo['unrealised_pnl']?.toString() ?? '0') ??
            0.0;
        final walletAssets = total + unrealizedPnL;
        return {'value': walletAssets, 'pnl': unrealizedPnL};
      }

      return {'value': 0.0, 'pnl': 0.0};
    } catch (e) {
      print('Error loading Gate.io Futures assets: $e');
      return {'value': 0.0, 'pnl': 0.0};
    }
  }

  Future<double> _loadBtcPrice() async {
    try {
      final priceData = await _binanceClient.getAllPrices();

      for (var price in priceData) {
        if (price['symbol'] == 'BTCUSDT') {
          return double.tryParse(price['price']) ?? 0.0;
        }
      }

      return 0.0;
    } catch (e) {
      print('Error loading BTC price: $e');
      return 0.0;
    }
  }

  @override
  Future<void> close() {
    _refreshTimer?.cancel();
    return super.close();
  }
}
