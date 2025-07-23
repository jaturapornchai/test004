class Asset {
  final String symbol;
  final String free;
  final String locked;
  final String exchange;
  final double usdValue;

  Asset({
    required this.symbol,
    required this.free,
    required this.locked,
    required this.exchange,
    required this.usdValue,
  });

  double get totalBalance => double.parse(free) + double.parse(locked);

  factory Asset.fromBinance(Map<String, dynamic> json, double usdPrice) {
    return Asset(
      symbol: json['asset'] ?? '',
      free: json['free'] ?? '0',
      locked: json['locked'] ?? '0',
      exchange: 'Binance',
      usdValue:
          (double.parse(json['free'] ?? '0') +
              double.parse(json['locked'] ?? '0')) *
          usdPrice,
    );
  }

  factory Asset.fromGate(Map<String, dynamic> json, double usdPrice) {
    return Asset(
      symbol: json['currency'] ?? '',
      free: json['available'] ?? '0',
      locked: json['locked'] ?? '0',
      exchange: 'Gate.io',
      usdValue:
          (double.parse(json['available'] ?? '0') +
              double.parse(json['locked'] ?? '0')) *
          usdPrice,
    );
  }
}
