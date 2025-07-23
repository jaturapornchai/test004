part of 'portfolio_cubit.dart';

enum PortfolioStatus { initial, loading, loaded, error }

class PortfolioState {
  final PortfolioStatus status;
  final double binanceSpotValue;
  final double binanceFuturesValue;
  final double binanceFuturesPnL;
  final double gateSpotValue;
  final double gateFuturesValue;
  final double gateFuturesPnL;
  final double usdToThbRate;
  final double btcPrice;
  final String? error;

  const PortfolioState({
    this.status = PortfolioStatus.initial,
    this.binanceSpotValue = 0.0,
    this.binanceFuturesValue = 0.0,
    this.binanceFuturesPnL = 0.0,
    this.gateSpotValue = 0.0,
    this.gateFuturesValue = 0.0,
    this.gateFuturesPnL = 0.0,
    this.usdToThbRate = 35.0,
    this.btcPrice = 0.0,
    this.error,
  });

  double get totalAmount =>
      binanceSpotValue + binanceFuturesValue + gateSpotValue + gateFuturesValue;

  double get totalPnL => binanceFuturesPnL + gateFuturesPnL;

  // Calculate how many BTC can be bought with total portfolio
  double get totalInBtc =>
      btcPrice > 0 ? (totalAmount + totalPnL) / btcPrice : 0.0;

  PortfolioState copyWith({
    PortfolioStatus? status,
    double? binanceSpotValue,
    double? binanceFuturesValue,
    double? binanceFuturesPnL,
    double? gateSpotValue,
    double? gateFuturesValue,
    double? gateFuturesPnL,
    double? usdToThbRate,
    double? btcPrice,
    String? error,
  }) {
    return PortfolioState(
      status: status ?? this.status,
      binanceSpotValue: binanceSpotValue ?? this.binanceSpotValue,
      binanceFuturesValue: binanceFuturesValue ?? this.binanceFuturesValue,
      binanceFuturesPnL: binanceFuturesPnL ?? this.binanceFuturesPnL,
      gateSpotValue: gateSpotValue ?? this.gateSpotValue,
      gateFuturesValue: gateFuturesValue ?? this.gateFuturesValue,
      gateFuturesPnL: gateFuturesPnL ?? this.gateFuturesPnL,
      usdToThbRate: usdToThbRate ?? this.usdToThbRate,
      btcPrice: btcPrice ?? this.btcPrice,
      error: error ?? this.error,
    );
  }
}
