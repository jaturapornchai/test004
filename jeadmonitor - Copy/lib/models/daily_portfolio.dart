class DailyPortfolio {
  final String date;
  final double totalValue;
  final double binanceSpot;
  final double binanceFutures;
  final double gateSpot;
  final double gateFutures;
  final double binanceFuturesPnL;
  final double gateFuturesPnL;

  DailyPortfolio({
    required this.date,
    required this.totalValue,
    required this.binanceSpot,
    required this.binanceFutures,
    required this.gateSpot,
    required this.gateFutures,
    required this.binanceFuturesPnL,
    required this.gateFuturesPnL,
  });

  Map<String, dynamic> toJson() {
    return {
      'date': date,
      'totalValue': totalValue,
      'binanceSpot': binanceSpot,
      'binanceFutures': binanceFutures,
      'gateSpot': gateSpot,
      'gateFutures': gateFutures,
      'binanceFuturesPnL': binanceFuturesPnL,
      'gateFuturesPnL': gateFuturesPnL,
    };
  }

  factory DailyPortfolio.fromJson(Map<String, dynamic> json) {
    return DailyPortfolio(
      date: json['date'] ?? '',
      totalValue: (json['totalValue'] ?? 0).toDouble(),
      binanceSpot: (json['binanceSpot'] ?? 0).toDouble(),
      binanceFutures: (json['binanceFutures'] ?? 0).toDouble(),
      gateSpot: (json['gateSpot'] ?? 0).toDouble(),
      gateFutures: (json['gateFutures'] ?? 0).toDouble(),
      binanceFuturesPnL: (json['binanceFuturesPnL'] ?? 0).toDouble(),
      gateFuturesPnL: (json['gateFuturesPnL'] ?? 0).toDouble(),
    );
  }
}
