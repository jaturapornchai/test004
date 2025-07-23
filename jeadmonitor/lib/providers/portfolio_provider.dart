import 'dart:async';
import 'package:flutter/foundation.dart';
import '../models/asset.dart';
import '../services/binance_api_client.dart';
import '../services/gate_api_client.dart';
import '../config/api_config.dart';

class PortfolioProvider with ChangeNotifier {
  List<Asset> _assets = [];
  bool _isLoading = false;
  String _error = '';
  Timer? _refreshTimer;

  late final BinanceApiClient _binanceClient;
  late final GateApiClient _gateClient;

  List<Asset> get assets => _assets;
  bool get isLoading => _isLoading;
  String get error => _error;

  double get totalUsdValue =>
      _assets.fold(0.0, (sum, asset) => sum + asset.usdValue);

  double get binanceTotalValue => _assets
      .where((asset) => asset.exchange == 'Binance')
      .fold(0.0, (sum, asset) => sum + asset.usdValue);

  double get gateTotalValue => _assets
      .where((asset) => asset.exchange == 'Gate.io')
      .fold(0.0, (sum, asset) => sum + asset.usdValue);

  // Spot และ Futures แยกกัน
  double get binanceSpotValue => _assets
      .where(
        (asset) =>
            asset.exchange == 'Binance' && asset.symbol.endsWith('_SPOT'),
      )
      .fold(0.0, (sum, asset) => sum + asset.usdValue);

  double get binanceFuturesValue => _assets
      .where(
        (asset) =>
            asset.exchange == 'Binance' && asset.symbol.endsWith('_FUTURES'),
      )
      .fold(0.0, (sum, asset) => sum + asset.usdValue);

  double get gateSpotValue => _assets
      .where(
        (asset) =>
            asset.exchange == 'Gate.io' && asset.symbol.endsWith('_SPOT'),
      )
      .fold(0.0, (sum, asset) => sum + asset.usdValue);

  double get gateFuturesValue => _assets
      .where(
        (asset) =>
            asset.exchange == 'Gate.io' && asset.symbol.endsWith('_FUTURES'),
      )
      .fold(0.0, (sum, asset) => sum + asset.usdValue);

  int get totalAssetsCount => _assets.length;

  int get binanceAssetsCount =>
      _assets.where((asset) => asset.exchange == 'Binance').length;

  int get gateAssetsCount =>
      _assets.where((asset) => asset.exchange == 'Gate.io').length;

  PortfolioProvider() {
    _initializeClients();
    // Remove auto-refresh since PortfolioScreen handles it completely
    // _startAutoRefresh();
  }

  void _initializeClients() {
    // Initialize with API configuration
    _binanceClient = BinanceApiClient(
      apiKey: ApiConfig.binanceApiKey,
      apiSecret: ApiConfig.binanceApiSecret,
    );

    _gateClient = GateApiClient(
      apiKey: ApiConfig.gateApiKey,
      apiSecret: ApiConfig.gateApiSecret,
    );
  }

  Future<void> refreshData() async {
    if (_isLoading) return;

    _isLoading = true;
    _error = '';
    notifyListeners();

    try {
      // โหลดข้อมูลจากทั้งสอง exchange และทั้ง Spot และ Futures
      final binanceSpotAssets = await _loadBinanceSpotAssets();
      final binanceFuturesAssets = await _loadBinanceFuturesAssets();
      final gateSpotAssets = await _loadGateSpotAssets();
      final gateFuturesAssets = await _loadGateFuturesAssets();

      _assets = [
        ...binanceSpotAssets,
        ...binanceFuturesAssets,
        ...gateSpotAssets,
        ...gateFuturesAssets,
      ];
      _assets.sort((a, b) => b.usdValue.compareTo(a.usdValue));
    } catch (e) {
      _error = e.toString();
      if (kDebugMode) {
        print('Error refreshing data: $e');
      }
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<List<Asset>> _loadBinanceSpotAssets() async {
    try {
      final accountInfo = await _binanceClient.getAccountInfo();
      final prices = await _binanceClient.getAllPrices();

      final priceMap = <String, double>{};
      for (final price in prices) {
        final symbol = price['symbol'] as String;
        if (symbol.endsWith('USDT')) {
          final baseSymbol = symbol.replaceAll('USDT', '');
          priceMap[baseSymbol] = double.tryParse(price['price'] ?? '0') ?? 0.0;
        }
      }
      priceMap['USDT'] = 1.0;

      final assets = <Asset>[];
      final balances = accountInfo['balances'] as List;

      for (final balance in balances) {
        final symbol = balance['asset'] as String;
        final free = balance['free'] ?? '0';
        final locked = balance['locked'] ?? '0';
        final totalBalance = double.parse(free) + double.parse(locked);

        if (totalBalance > 0) {
          final usdPrice = priceMap[symbol] ?? 0.0;
          final usdValue = totalBalance * usdPrice;

          if (usdValue > 0.01) {
            assets.add(
              Asset(
                symbol: '${symbol}_SPOT',
                free: free,
                locked: locked,
                exchange: 'Binance',
                usdValue: usdValue,
              ),
            );
          }
        }
      }

      return assets;
    } catch (e) {
      if (kDebugMode) {
        print('Error loading Binance Spot assets: $e');
      }
      return [];
    }
  }

  Future<List<Asset>> _loadBinanceFuturesAssets() async {
    try {
      final response = await _binanceClient.getFuturesAccount();
      final assets = <Asset>[];

      if (response != null) {
        final assetsList = response['assets'] as List;

        for (final asset in assetsList) {
          final symbol = asset['asset'] as String;
          final balance = double.tryParse(asset['walletBalance'] ?? '0') ?? 0.0;
          final pnl = double.tryParse(asset['unrealizedProfit'] ?? '0') ?? 0.0;

          if (balance > 0.01) {
            assets.add(
              Asset(
                symbol: '${symbol}_FUTURES',
                free: balance.toString(),
                locked: '0',
                exchange: 'Binance',
                usdValue: balance + pnl,
              ),
            );
          }
        }
      }

      return assets;
    } catch (e) {
      if (kDebugMode) {
        print('Error loading Binance Futures assets: $e');
      }
      return [];
    }
  }

  Future<List<Asset>> _loadGateSpotAssets() async {
    try {
      final accounts = await _gateClient.getSpotAccounts();
      final assets = <Asset>[];

      for (final account in accounts) {
        final currency = account['currency'] as String;
        final available = account['available'] ?? '0';
        final locked = account['locked'] ?? '0';
        final total = double.parse(available) + double.parse(locked);

        if (total > 0) {
          double usdValue = total;
          if (currency == 'USDT' || currency == 'USD') {
            usdValue = total;
          } else {
            usdValue = total * 1.0; // ควรใช้ราคาจริง
          }

          if (usdValue > 0.01) {
            assets.add(
              Asset(
                symbol: '${currency}_SPOT',
                free: available,
                locked: locked,
                exchange: 'Gate.io',
                usdValue: usdValue,
              ),
            );
          }
        }
      }

      return assets;
    } catch (e) {
      if (kDebugMode) {
        print('Error loading Gate.io Spot assets: $e');
      }
      return [];
    }
  }

  Future<List<Asset>> _loadGateFuturesAssets() async {
    try {
      final response = await _gateClient.getFuturesAccount();
      final assets = <Asset>[];

      if (response != null) {
        final total =
            double.tryParse(response['total']?.toString() ?? '0') ?? 0.0;
        final unrealizedPnl =
            double.tryParse(response['unrealised_pnl']?.toString() ?? '0') ??
            0.0;

        if (total > 0.01) {
          assets.add(
            Asset(
              symbol: 'USDT_FUTURES',
              free: total.toString(),
              locked: '0',
              exchange: 'Gate.io',
              usdValue: total + unrealizedPnl,
            ),
          );
        }
      }

      return assets;
    } catch (e) {
      if (kDebugMode) {
        print('Error loading Gate.io Futures assets: $e');
      }
      return [];
    }
  }

  @override
  void dispose() {
    _refreshTimer?.cancel();
    super.dispose();
  }
}
