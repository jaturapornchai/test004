import 'dart:convert';
import 'package:http/http.dart' as http;

class PriceService {
  static Map<String, double> _cachedPrices = {};
  static DateTime? _lastUpdate;

  static Future<Map<String, double>> getCurrentPrices() async {
    // Cache prices for 5 minutes
    if (_lastUpdate != null &&
        DateTime.now().difference(_lastUpdate!).inMinutes < 5) {
      return _cachedPrices;
    }

    try {
      const url =
          'https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum,cardano,binancecoin,usd-coin,tether,solana&vs_currencies=usd';
      final response = await http.get(Uri.parse(url));

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body) as Map<String, dynamic>;
        _cachedPrices = {
          'BTC': (data['bitcoin']?['usd'] ?? 45000).toDouble(),
          'ETH': (data['ethereum']?['usd'] ?? 2500).toDouble(),
          'ADA': (data['cardano']?['usd'] ?? 0.4).toDouble(),
          'BNB': (data['binancecoin']?['usd'] ?? 300).toDouble(),
          'USDC': (data['usd-coin']?['usd'] ?? 1.0).toDouble(),
          'USDT': (data['tether']?['usd'] ?? 1.0).toDouble(),
          'SOL': (data['solana']?['usd'] ?? 100).toDouble(),
        };
        _lastUpdate = DateTime.now();
        return _cachedPrices;
      }
    } catch (e) {
      print('Failed to fetch prices: $e');
    }

    // Fallback prices
    return {
      'BTC': 45000.0,
      'ETH': 2500.0,
      'ADA': 0.4,
      'BNB': 300.0,
      'USDC': 1.0,
      'USDT': 1.0,
      'SOL': 100.0,
    };
  }

  static double getPrice(String symbol) {
    return _cachedPrices[symbol] ?? 1.0;
  }
}
