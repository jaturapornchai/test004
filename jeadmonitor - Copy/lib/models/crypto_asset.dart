class CryptoAsset {
  final String symbol;
  final double amount;
  final double usdValue;
  final String exchange;
  final String marketType;

  CryptoAsset({
    required this.symbol,
    required this.amount,
    required this.usdValue,
    required this.exchange,
    required this.marketType,
  });

  factory CryptoAsset.fromJson(Map<String, dynamic> json) {
    return CryptoAsset(
      symbol: json['symbol'] ?? '',
      amount: double.tryParse(json['amount'].toString()) ?? 0.0,
      usdValue: double.tryParse(json['usdValue'].toString()) ?? 0.0,
      exchange: json['exchange'] ?? '',
      marketType: json['marketType'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'symbol': symbol,
      'amount': amount,
      'usdValue': usdValue,
      'exchange': exchange,
      'marketType': marketType,
    };
  }

  @override
  String toString() {
    return 'CryptoAsset(symbol: $symbol, amount: $amount, usdValue: $usdValue, exchange: $exchange, marketType: $marketType)';
  }
}
