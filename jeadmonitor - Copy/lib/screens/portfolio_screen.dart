import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:wakelock_plus/wakelock_plus.dart';
import 'package:intl/intl.dart';
import '../models/crypto_asset.dart';
import '../services/binance_api_client.dart';
import '../services/gate_api_client.dart';
import '../services/exchange_rate_service.dart';
import '../utils/text_style_helper.dart';

class PortfolioScreen extends StatefulWidget {
  const PortfolioScreen({Key? key}) : super(key: key);

  @override
  State<PortfolioScreen> createState() => _PortfolioScreenState();
}

class _PortfolioScreenState extends State<PortfolioScreen> {
  final BinanceApiClient _binanceClient = BinanceApiClient(
    apiKey: 'xz3oxKAcLPRJxgQvBOSqNGyJipzOoiAvRXwvaHZBYSyHFEfegRwV37Blmwwpxp4z',
    apiSecret:
        'h9LCIdryqVCEPn6JgGSnmycgtc7FRtKnrSVcqOFVmJ3bOChZDvGBcOCnCdOhvFlT',
  );
  final GateApiClient _gateClient = GateApiClient(
    apiKey: '472667d2eebf07b206f8599bf35bbcb1',
    apiSecret:
        'e24966020b7241e521966b87deb401d6dac0f90a1ccbf3f85db54ecc09faf33f',
  );
  final ExchangeRateService _exchangeRateService = ExchangeRateService();

  List<CryptoAsset> _binanceSpotAssets = [];
  List<CryptoAsset> _binanceFuturesAssets = [];
  List<CryptoAsset> _gateSpotAssets = [];
  // Gate.io Futures data
  double _gateFuturesTotal = 0.0;

  // PnL data
  double _binanceFuturesPnL = 0.0;
  double _gateFuturesPnL = 0.0;

  // Exchange rate
  double _usdToThbRate = 35.0; // Default rate

  Timer? _refreshTimer;
  Timer? _exchangeRateTimer;

  @override
  void initState() {
    super.initState();

    // Set landscape orientation and fullscreen
    SystemChrome.setPreferredOrientations([
      DeviceOrientation.landscapeLeft,
      DeviceOrientation.landscapeRight,
    ]);
    SystemChrome.setEnabledSystemUIMode(SystemUiMode.immersiveSticky);

    // Keep screen on and set maximum brightness
    WakelockPlus.enable();
    SystemChrome.setSystemUIOverlayStyle(
      const SystemUiOverlayStyle(statusBarBrightness: Brightness.light),
    );

    _loadData(); // Initial load with loading indicator
    _loadExchangeRate(); // Load initial exchange rate

    // Auto-refresh every 5 seconds without loading indicator
    _refreshTimer = Timer.periodic(const Duration(seconds: 5), (timer) {
      _loadData(showLoading: false); // Refresh without showing loading
    });

    // Update exchange rate every 1 hour
    _exchangeRateTimer = Timer.periodic(const Duration(hours: 1), (timer) {
      _loadExchangeRate();
    });
  }

  @override
  void dispose() {
    _refreshTimer?.cancel();
    _exchangeRateTimer?.cancel();
    WakelockPlus.disable();
    super.dispose();
  }

  Future<void> _loadData({bool showLoading = true}) async {
    try {
      // Load all data in parallel
      await Future.wait([
        _loadBinanceSpotAssets(),
        _loadBinanceFuturesAssets(),
        _loadGateSpotAssets(),
        _loadGateFuturesAssets(),
      ]);
    } catch (e) {
      print('Failed to load portfolio data: $e');
    }
  }

  Future<void> _loadExchangeRate() async {
    try {
      final rate = await _exchangeRateService.getUsdToThbRate();
      setState(() {
        _usdToThbRate = rate;
      });
      print('Exchange rate updated: 1 USD = $_usdToThbRate THB');
    } catch (e) {
      print('Failed to load exchange rate: $e');
    }
  }

  Future<void> _loadBinanceSpotAssets() async {
    try {
      final accountInfo = await _binanceClient.getAccountInfo();
      final List<CryptoAsset> assets = [];

      // Get all prices first
      Map<String, double> prices = {};
      try {
        final priceData = await _binanceClient.getAllPrices();
        for (var price in priceData) {
          final symbol = price['symbol'] as String;
          final priceValue = double.tryParse(price['price']) ?? 0.0;
          if (symbol.endsWith('USDT')) {
            final asset = symbol.replaceAll('USDT', '');
            prices[asset] = priceValue;
          }
        }
      } catch (e) {
        print('Error loading Binance prices: $e');
      }

      // Calculate total spot value in USDT
      double totalSpotValue = 0.0;
      for (var balance in accountInfo['balances']) {
        final double free = double.tryParse(balance['free']) ?? 0.0;
        final double locked = double.tryParse(balance['locked']) ?? 0.0;
        final double total = free + locked;
        final String asset = balance['asset'];

        if (total > 0) {
          double usdValue = 0.0;
          if (asset == 'USDT' || asset == 'USDC' || asset == 'BUSD') {
            // Stablecoins are 1:1 with USD
            usdValue = total;
          } else if (prices.containsKey(asset)) {
            // Calculate USD value using price
            usdValue = total * prices[asset]!;
          }

          if (usdValue > 0.01) {
            // Only count if value > 1 cent
            totalSpotValue += usdValue;
          }
        }
      }

      // Add only one entry representing total spot portfolio
      if (totalSpotValue > 0) {
        assets.add(
          CryptoAsset(
            symbol: 'SPOT',
            amount: 1.0, // Representing total spot portfolio
            usdValue: totalSpotValue, // Use actual USDT total
            exchange: 'Binance',
            marketType: 'Spot',
          ),
        );
      }

      setState(() {
        _binanceSpotAssets = assets;
      });
    } catch (e) {
      print('Error loading Binance Spot assets: $e');
    }
  }

  Future<void> _loadBinanceFuturesAssets() async {
    try {
      final accountInfo = await _binanceClient.getFuturesAccount();
      final List<CryptoAsset> assets = [];

      // Get balance and PnL separately (like in screenshot)
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

      // Calculate balance WITHOUT PnL (like screenshot shows Balance: 197.8223)
      final balanceWithoutPnL = totalWalletBalance - totalUnrealizedProfit;

      _binanceFuturesPnL = totalUnrealizedProfit;

      if (balanceWithoutPnL > 0) {
        assets.add(
          CryptoAsset(
            symbol: 'FUTURES',
            amount: 1.0,
            usdValue: balanceWithoutPnL + _binanceFuturesPnL,
            exchange: 'Binance',
            marketType: 'Futures',
          ),
        );
      }

      setState(() {
        _binanceFuturesAssets = assets;
      });
    } catch (e) {
      print('Error loading Binance Futures assets: $e');
    }
  }

  Future<void> _loadGateSpotAssets() async {
    try {
      final accounts = await _gateClient.getSpotAccounts();
      final List<CryptoAsset> assets = [];

      // Get all prices first
      Map<String, double> prices = {};
      try {
        final tickers = await _gateClient.getAllTickers();
        for (var ticker in tickers) {
          final symbol = ticker['currency_pair'] as String;
          final lastPrice =
              double.tryParse(ticker['last']?.toString() ?? '0') ?? 0.0;
          if (symbol.endsWith('_USDT')) {
            final asset = symbol.replaceAll('_USDT', '');
            prices[asset] = lastPrice;
          }
        }
      } catch (e) {
        print('Error loading Gate.io prices: $e');
      }

      // Calculate total spot value in USDT
      double totalSpotValue = 0.0;
      for (var account in accounts) {
        final double available = double.tryParse(account['available']) ?? 0.0;
        final double locked = double.tryParse(account['locked']) ?? 0.0;
        final double total = available + locked;
        final String currency = account['currency'] ?? '';

        if (total > 0) {
          double usdValue = 0.0;
          if (currency == 'USDT' || currency == 'USDC') {
            // Stablecoins are 1:1 with USD
            usdValue = total;
          } else if (prices.containsKey(currency)) {
            // Calculate USD value using price
            usdValue = total * prices[currency]!;
          }

          if (usdValue > 0.01) {
            // Only count if value > 1 cent
            totalSpotValue += usdValue;
          }
        }
      }

      // Add only one entry representing total spot portfolio
      if (totalSpotValue > 0) {
        assets.add(
          CryptoAsset(
            symbol: 'SPOT',
            amount: 1.0, // Representing total spot portfolio
            usdValue: totalSpotValue, // Use actual USDT total
            exchange: 'Gate.io',
            marketType: 'Spot',
          ),
        );
      }

      setState(() {
        _gateSpotAssets = assets;
      });
    } catch (e) {
      print('Error loading Gate.io Spot assets: $e');
    }
  }

  Future<void> _loadGateFuturesAssets() async {
    try {
      final accountInfo = await _gateClient.getFuturesAccount();

      if (accountInfo != null) {
        final walletAssets =
            double.tryParse(accountInfo['total']?.toString() ?? '0') ?? 0.0;
        final unrealizedPnL =
            double.tryParse(accountInfo['unrealised_pnl']?.toString() ?? '0') ??
            0.0;

        // For consistency with screenshot: use wallet assets as futures value
        // Total will be calculated as spot + futures + pnl separately
        setState(() {
          _gateFuturesTotal = walletAssets; // Wallet Assets (1206.06)
          _gateFuturesPnL = unrealizedPnL; // PnL (+29.24)
        });
      }
    } catch (e) {
      print('Error loading Gate.io Futures assets: $e');
    }
  }

  double get _binanceSpotValue =>
      _binanceSpotAssets.fold(0.0, (sum, asset) => sum + asset.usdValue);
  double get _binanceFuturesValue =>
      _binanceFuturesAssets.fold(0.0, (sum, asset) => sum + asset.usdValue);
  double get _gateSpotValue =>
      _gateSpotAssets.fold(0.0, (sum, asset) => sum + asset.usdValue);
  double get _gateFuturesValue => _gateFuturesTotal;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF0A0E1A), Color(0xFF1A1A2E), Color(0xFF16213E)],
          ),
        ),
        child: SafeArea(
          child: Padding(
            padding: const EdgeInsets.all(2.0),
            child: Row(
              children: [
                // Column 1: ยอดเงินรวม
                Expanded(flex: 1, child: _buildTotalSummaryCard()),
                const SizedBox(width: 4),
                // Column 2: Binance
                Expanded(
                  flex: 1,
                  child: _buildDetailedExchangeCard(
                    'Binance',
                    _binanceSpotValue,
                    _binanceFuturesValue,
                    Colors.orange,
                    _binanceFuturesPnL,
                  ),
                ),
                const SizedBox(width: 4),
                // Column 3: Gate.io
                Expanded(
                  flex: 1,
                  child: _buildDetailedExchangeCard(
                    'Gate.io',
                    _gateSpotValue,
                    _gateFuturesValue,
                    Colors.green,
                    _gateFuturesPnL,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildTotalSummaryCard() {
    return Container(
      width: double.infinity,
      height: double.infinity,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF1E3A8A), Color(0xFF1D4ED8)],
        ),
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.blue.withOpacity(0.2),
            blurRadius: 8,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          _buildTotalAmountText(),
          const SizedBox(height: 38),
          _buildPnLText(),
        ],
      ),
    );
  }

  Widget _buildTotalAmountText() {
    final totalAmount =
        _binanceSpotValue +
        _binanceFuturesValue +
        _gateSpotValue +
        _gateFuturesValue +
        _binanceFuturesPnL +
        _gateFuturesPnL;
    final totalAmountThb = totalAmount * _usdToThbRate;
    final formatter = NumberFormat('#,##0.00');

    return Column(
      children: [
        Text(
          'ยอดเงินรวม',
          style: TextStyleHelper.withOutlineAndShadow(
            fontSize: 22,
            fontWeight: FontWeight.bold,
            color: Colors.white,
            outlineWidth: 0.8,
            shadowBlurRadius: 3.0,
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 4),
        Text(
          '\$${formatter.format(totalAmount)}',
          style: TextStyleHelper.withOutlineAndShadow(
            fontSize: 32,
            fontWeight: FontWeight.bold,
            color: Colors.white,
            outlineWidth: 0.8,
            shadowBlurRadius: 3.0,
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 4),
        Text(
          '฿${formatter.format(totalAmountThb)}',
          style: TextStyleHelper.withOutlineAndShadow(
            fontSize: 24,
            fontWeight: FontWeight.w600,
            color: Colors.white70,
            outlineWidth: 0.6,
            shadowBlurRadius: 2.5,
          ),
          textAlign: TextAlign.center,
        ),
      ],
    );
  }

  Widget _buildPnLText() {
    final unrealizedPnL = _binanceFuturesPnL + _gateFuturesPnL;
    final unrealizedPnLThb = unrealizedPnL * _usdToThbRate;
    final formatter = NumberFormat('#,##0.00');

    // Calculate percentage if possible
    final totalPortfolioValue =
        _binanceSpotValue +
        _binanceFuturesValue +
        _gateSpotValue +
        _gateFuturesValue;
    final pnlPercentage = totalPortfolioValue > 0
        ? (unrealizedPnL / totalPortfolioValue) * 100
        : 0.0;

    return Column(
      children: [
        // Line 1: Title
        Text(
          'กำไร/ขาดทุนปัจจุบัน',
          style: TextStyleHelper.withOutlineAndShadow(
            fontSize: 20,
            fontWeight: FontWeight.bold,
            color: Colors.white,
            outlineWidth: 0.8,
            shadowBlurRadius: 3.0,
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 6),
        // Line 2: USD Amount
        Text(
          '\$${formatter.format(unrealizedPnL)}',
          style: TextStyleHelper.withOutlineAndShadow(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: unrealizedPnL >= 0 ? Colors.green : Colors.red,
            outlineWidth: 0.9,
            shadowBlurRadius: 3.5,
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 4),
        // Line 3: THB Amount and Percentage
        Text(
          '฿${formatter.format(unrealizedPnLThb)} (${pnlPercentage >= 0 ? '+' : ''}${pnlPercentage.toStringAsFixed(2)}%)',
          style: TextStyleHelper.withOutlineAndShadow(
            fontSize: 18,
            fontWeight: FontWeight.w600,
            color: unrealizedPnL >= 0
                ? Colors.green.withOpacity(0.9)
                : Colors.red.withOpacity(0.9),
            outlineWidth: 0.7,
            shadowBlurRadius: 2.5,
          ),
          textAlign: TextAlign.center,
        ),
      ],
    );
  }

  Widget _buildDetailedExchangeCard(
    String exchange,
    double spotValue,
    double futuresValue,
    Color accentColor,
    double futuresPnL,
  ) {
    return Container(
      width: double.infinity,
      height: double.infinity,
      padding: const EdgeInsets.all(6),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            accentColor.withOpacity(0.15),
            accentColor.withOpacity(0.05),
          ],
        ),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: accentColor.withOpacity(0.4), width: 2),
        boxShadow: [
          BoxShadow(
            color: accentColor.withOpacity(0.2),
            blurRadius: 8,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          _buildExchangeTitle(exchange, accentColor),
          const SizedBox(height: 2),
          Expanded(
            child: _buildExchangeDetails(
              spotValue,
              futuresValue,
              futuresPnL,
              accentColor,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildExchangeTitle(String exchange, Color accentColor) {
    return Text(
      exchange,
      style: TextStyleHelper.withOutlineAndShadow(
        fontSize: 7.5,
        fontWeight: FontWeight.bold,
        color: accentColor,
        outlineWidth: 0.8,
        shadowBlurRadius: 3.0,
      ),
      textAlign: TextAlign.center,
    );
  }

  Widget _buildExchangeDetails(
    double spotValue,
    double futuresValue,
    double futuresPnL,
    Color accentColor,
  ) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.1),
        borderRadius: BorderRadius.circular(10),
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          _buildValueRow('Spot:', spotValue),
          _buildValueRow('Futures:', futuresValue),
          _buildPnLRow(futuresPnL),
          const Divider(color: Colors.white30, height: 12),
          _buildTotalRow(spotValue + futuresValue + futuresPnL, accentColor),
        ],
      ),
    );
  }

  Widget _buildValueRow(String label, double value) {
    final valueThb = value * _usdToThbRate;
    final formatter = NumberFormat('#,##0.00');
    final formatterThb = NumberFormat('#,##0');

    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              label,
              style: TextStyleHelper.withOutlineAndShadow(
                fontSize: 14.5,
                fontWeight: FontWeight.w600,
                color: Colors.white70,
                outlineWidth: 0.4,
                shadowBlurRadius: 2.0,
              ),
            ),
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  '\$${formatter.format(value)}',
                  style: TextStyleHelper.withOutlineAndShadow(
                    fontSize: 22.5,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                    outlineWidth: 0.6,
                    shadowBlurRadius: 2.5,
                  ),
                ),
                Text(
                  '฿${formatterThb.format(valueThb)}',
                  style: TextStyleHelper.withOutlineAndShadow(
                    fontSize: 18,
                    fontWeight: FontWeight.w500,
                    color: Colors.white60,
                    outlineWidth: 0.3,
                    shadowBlurRadius: 1.5,
                  ),
                ),
              ],
            ),
          ],
        ),
        const SizedBox(height: 4),
      ],
    );
  }

  Widget _buildPnLRow(double pnl) {
    final pnlThb = pnl * _usdToThbRate;
    final formatter = NumberFormat('#,##0.00');
    final formatterThb = NumberFormat('#,##0');

    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              'PnL:',
              style: TextStyleHelper.withOutlineAndShadow(
                fontSize: 14.5,
                fontWeight: FontWeight.w600,
                color: Colors.white70,
                outlineWidth: 0.3,
                shadowBlurRadius: 1.5,
              ),
            ),
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  '${pnl >= 0 ? '+' : ''}\$${formatter.format(pnl.abs())}',
                  style: TextStyleHelper.withOutlineAndShadow(
                    fontSize: 19.5,
                    fontWeight: FontWeight.bold,
                    color: pnl >= 0 ? Colors.green : Colors.red,
                    outlineWidth: 0.5,
                    shadowBlurRadius: 2.0,
                  ),
                ),
                Text(
                  '${pnl >= 0 ? '+' : ''}฿${formatterThb.format(pnlThb.abs())}',
                  style: TextStyleHelper.withOutlineAndShadow(
                    fontSize: 16.5,
                    fontWeight: FontWeight.w500,
                    color: pnl >= 0
                        ? Colors.green.withOpacity(0.8)
                        : Colors.red.withOpacity(0.8),
                    outlineWidth: 0.3,
                    shadowBlurRadius: 1.5,
                  ),
                ),
              ],
            ),
          ],
        ),
        const SizedBox(height: 4),
      ],
    );
  }

  Widget _buildTotalRow(double total, Color accentColor) {
    final totalThb = total * _usdToThbRate;
    final formatter = NumberFormat('#,##0.00');
    final formatterThb = NumberFormat('#,##0');

    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Text(
          '\$${formatter.format(total)}',
          style: TextStyleHelper.withOutlineAndShadow(
            fontSize: 25.5,
            fontWeight: FontWeight.bold,
            color: accentColor,
            outlineWidth: 0.7,
            shadowBlurRadius: 3.0,
          ),
        ),
        Text(
          '฿${formatterThb.format(totalThb)}',
          style: TextStyleHelper.withOutlineAndShadow(
            fontSize: 19.5,
            fontWeight: FontWeight.w600,
            color: accentColor.withOpacity(0.8),
            outlineWidth: 0.4,
            shadowBlurRadius: 2.0,
          ),
        ),
      ],
    );
  }
}
