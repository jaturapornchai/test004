import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:wakelock_plus/wakelock_plus.dart';
import 'package:intl/intl.dart';
import '../cubits/portfolio_cubit.dart';
import '../utils/text_style_helper.dart';

class PortfolioScreenCubit extends StatefulWidget {
  const PortfolioScreenCubit({Key? key}) : super(key: key);

  @override
  State<PortfolioScreenCubit> createState() => _PortfolioScreenCubitState();
}

class _PortfolioScreenCubitState extends State<PortfolioScreenCubit> {
  @override
  void initState() {
    super.initState();

    // Set landscape orientation and fullscreen
    SystemChrome.setPreferredOrientations([
      DeviceOrientation.landscapeLeft,
      DeviceOrientation.landscapeRight,
    ]);
    SystemChrome.setEnabledSystemUIMode(SystemUiMode.immersiveSticky);

    // Keep screen on
    WakelockPlus.enable();
    SystemChrome.setSystemUIOverlayStyle(
      const SystemUiOverlayStyle(statusBarBrightness: Brightness.light),
    );

    // Initialize portfolio data
    context.read<PortfolioCubit>().initialize();
  }

  @override
  void dispose() {
    WakelockPlus.disable();
    super.dispose();
  }

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
            child: BlocBuilder<PortfolioCubit, PortfolioState>(
              builder: (context, state) {
                if (state.status == PortfolioStatus.loading) {
                  return const Center(
                    child: CircularProgressIndicator(color: Colors.blue),
                  );
                }

                if (state.status == PortfolioStatus.error) {
                  return Center(
                    child: Text(
                      'Error: ${state.error}',
                      style: const TextStyle(color: Colors.red, fontSize: 16),
                      textAlign: TextAlign.center,
                    ),
                  );
                }

                return Row(
                  children: [
                    // Column 1: ยอดเงินรวม
                    Expanded(flex: 1, child: _buildTotalSummaryCard(state)),
                    const SizedBox(width: 4),
                    // Column 2: Binance
                    Expanded(
                      flex: 1,
                      child: _buildDetailedExchangeCard(
                        'Binance',
                        state.binanceSpotValue,
                        state.binanceFuturesValue,
                        Colors.orange,
                        state.binanceFuturesPnL,
                        state.usdToThbRate,
                      ),
                    ),
                    const SizedBox(width: 4),
                    // Column 3: Gate.io
                    Expanded(
                      flex: 1,
                      child: _buildDetailedExchangeCard(
                        'Gate.io',
                        state.gateSpotValue,
                        state.gateFuturesValue,
                        Colors.green,
                        state.gateFuturesPnL,
                        state.usdToThbRate,
                      ),
                    ),
                  ],
                );
              },
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildTotalSummaryCard(PortfolioState state) {
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
          _buildTotalAmountText(state),
          const SizedBox(height: 18),
          _buildPnLText(state),
        ],
      ),
    );
  }

  Widget _buildTotalAmountText(PortfolioState state) {
    final totalAmountThb = state.totalAmount * state.usdToThbRate;
    final formatter = NumberFormat('#,##0.00');

    return Column(
      children: [
        Text(
          '\$${formatter.format(state.totalAmount)}',
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
        const SizedBox(height: 8),
        _buildBtcEquivalent(state),
      ],
    );
  }

  Widget _buildPnLText(PortfolioState state) {
    final unrealizedPnLThb = state.totalPnL * state.usdToThbRate;
    final formatter = NumberFormat('#,##0.00');

    // Calculate percentage
    final totalPortfolioValue =
        state.binanceSpotValue +
        state.binanceFuturesValue +
        state.gateSpotValue +
        state.gateFuturesValue;
    final pnlPercentage = totalPortfolioValue > 0
        ? (state.totalPnL / totalPortfolioValue) * 100
        : 0.0;

    return Column(
      children: [
        Text(
          'กำไร/ขาดทุน',
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
        Text(
          '\$${formatter.format(state.totalPnL)}',
          style: TextStyleHelper.withOutlineAndShadow(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: state.totalPnL >= 0 ? Colors.green : Colors.red,
            outlineWidth: 0.9,
            shadowBlurRadius: 3.5,
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 4),
        Text(
          '฿${formatter.format(unrealizedPnLThb)} (${pnlPercentage >= 0 ? '+' : ''}${pnlPercentage.toStringAsFixed(2)}%)',
          style: TextStyleHelper.withOutlineAndShadow(
            fontSize: 18,
            fontWeight: FontWeight.w600,
            color: state.totalPnL >= 0
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
    double usdToThbRate,
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
              usdToThbRate,
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
    double usdToThbRate,
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
          _buildValueRow('Spot:', spotValue, usdToThbRate),
          _buildValueRow('Futures:', futuresValue, usdToThbRate),
          _buildPnLRow(futuresPnL, usdToThbRate),
          const Divider(color: Colors.white30, height: 12),
          _buildTotalRow((spotValue + futuresValue), accentColor, usdToThbRate),
        ],
      ),
    );
  }

  Widget _buildValueRow(String label, double value, double usdToThbRate) {
    final valueThb = value * usdToThbRate;
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

  Widget _buildPnLRow(double pnl, double usdToThbRate) {
    final pnlThb = pnl * usdToThbRate;
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

  Widget _buildTotalRow(double total, Color accentColor, double usdToThbRate) {
    final totalThb = total * usdToThbRate;
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

  Widget _buildBtcEquivalent(PortfolioState state) {
    final formatter = NumberFormat('#,##0.00000000');

    return Column(
      children: [
        Text(
          'BTC: \$${NumberFormat('#,##0.00').format(state.btcPrice)}',
          style: TextStyleHelper.withOutlineAndShadow(
            fontSize: 16,
            fontWeight: FontWeight.w500,
            color: Colors.orange.shade300,
            outlineWidth: 0.5,
            shadowBlurRadius: 2.0,
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 2),
        Text(
          '≈ ${formatter.format(state.totalInBtc)} BTC',
          style: TextStyleHelper.withOutlineAndShadow(
            fontSize: 18,
            fontWeight: FontWeight.w600,
            color: Colors.orange.shade200,
            outlineWidth: 0.6,
            shadowBlurRadius: 2.5,
          ),
          textAlign: TextAlign.center,
        ),
      ],
    );
  }
}
