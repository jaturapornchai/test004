import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'package:http/http.dart' as http;

class BinanceApiClient {
  static const String baseUrl = 'https://api.binance.com';
  static const String futuresUrl = 'https://fapi.binance.com';
  final String apiKey;
  final String apiSecret;

  BinanceApiClient({required this.apiKey, required this.apiSecret});

  String _createSignature(String query) {
    var key = utf8.encode(apiSecret);
    var bytes = utf8.encode(query);
    var hmacSha256 = Hmac(sha256, key);
    var digest = hmacSha256.convert(bytes);
    return digest.toString();
  }

  // Public method to generate signature (for external use)
  String generateSignature(String query) {
    return _createSignature(query);
  }

  Future<Map<String, dynamic>> getAccountInfo() async {
    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch.toString();
      final query = 'timestamp=$timestamp&recvWindow=5000';
      final signature = _createSignature(query);

      final response = await http
          .get(
            Uri.parse('$baseUrl/api/v3/account?$query&signature=$signature'),
            headers: {
              'X-MBX-APIKEY': apiKey,
              'Content-Type': 'application/json',
            },
          )
          .timeout(const Duration(seconds: 15));

      if (response.statusCode == 200) {
        return json.decode(response.body);
      } else {
        print('Binance API Error: ${response.statusCode} - ${response.body}');
        throw Exception(
          'Failed to load account info: ${response.statusCode} - ${response.body}',
        );
      }
    } catch (e) {
      print('Binance Network Error: $e');
      throw Exception('Network error: $e');
    }
  }

  Future<Map<String, dynamic>?> getFuturesAccount() async {
    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch.toString();
      final query = 'timestamp=$timestamp&recvWindow=5000';
      final signature = _createSignature(query);

      final response = await http
          .get(
            Uri.parse(
              '$futuresUrl/fapi/v2/account?$query&signature=$signature',
            ),
            headers: {
              'X-MBX-APIKEY': apiKey,
              'Content-Type': 'application/json',
            },
          )
          .timeout(const Duration(seconds: 15));

      if (response.statusCode == 200) {
        return json.decode(response.body);
      } else {
        print(
          'Binance Futures API Error: ${response.statusCode} - ${response.body}',
        );
        return null;
      }
    } catch (e) {
      print('Binance Futures Network Error: $e');
      return null;
    }
  }

  Future<List<Map<String, dynamic>>> getAllPrices() async {
    try {
      final response = await http
          .get(Uri.parse('$baseUrl/api/v3/ticker/price'))
          .timeout(const Duration(seconds: 10));

      if (response.statusCode == 200) {
        return List<Map<String, dynamic>>.from(json.decode(response.body));
      } else {
        throw Exception('Failed to load prices: ${response.statusCode}');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  // Get 24hr ticker statistics for price change data
  Future<List<Map<String, dynamic>>> get24hrTicker() async {
    try {
      final response = await http
          .get(Uri.parse('$baseUrl/api/v3/ticker/24hr'))
          .timeout(const Duration(seconds: 10));

      if (response.statusCode == 200) {
        return List<Map<String, dynamic>>.from(json.decode(response.body));
      } else {
        throw Exception('Failed to load 24hr ticker: ${response.statusCode}');
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  // Get 24hr ticker for futures
  Future<List<Map<String, dynamic>>> getFutures24hrTicker() async {
    try {
      final response = await http
          .get(Uri.parse('$futuresUrl/fapi/v1/ticker/24hr'))
          .timeout(const Duration(seconds: 10));

      if (response.statusCode == 200) {
        return List<Map<String, dynamic>>.from(json.decode(response.body));
      } else {
        throw Exception(
          'Failed to load futures 24hr ticker: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  // Get today's futures trade history (closed positions)
  Future<List<Map<String, dynamic>>> getTodayFuturesTradeHistory() async {
    try {
      final now = DateTime.now();
      final todayStart = DateTime(now.year, now.month, now.day);
      final startTime = todayStart.millisecondsSinceEpoch;
      final endTime = now.millisecondsSinceEpoch;

      final timestamp = DateTime.now().millisecondsSinceEpoch.toString();
      final query =
          'startTime=$startTime&endTime=$endTime&timestamp=$timestamp&recvWindow=5000';
      final signature = _createSignature(query);

      final response = await http
          .get(
            Uri.parse(
              '$futuresUrl/fapi/v1/userTrades?$query&signature=$signature',
            ),
            headers: {
              'X-MBX-APIKEY': apiKey,
              'Content-Type': 'application/json',
            },
          )
          .timeout(const Duration(seconds: 15));

      if (response.statusCode == 200) {
        return List<Map<String, dynamic>>.from(json.decode(response.body));
      } else {
        throw Exception(
          'Failed to load futures trade history: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }

  // Get today's spot trade history
  Future<List<Map<String, dynamic>>> getTodaySpotTradeHistory() async {
    try {
      final now = DateTime.now();
      final todayStart = DateTime(now.year, now.month, now.day);
      final startTime = todayStart.millisecondsSinceEpoch;
      final endTime = now.millisecondsSinceEpoch;

      final timestamp = DateTime.now().millisecondsSinceEpoch.toString();
      final query =
          'startTime=$startTime&endTime=$endTime&timestamp=$timestamp&recvWindow=5000';
      final signature = _createSignature(query);

      final response = await http
          .get(
            Uri.parse('$baseUrl/api/v3/myTrades?$query&signature=$signature'),
            headers: {
              'X-MBX-APIKEY': apiKey,
              'Content-Type': 'application/json',
            },
          )
          .timeout(const Duration(seconds: 15));

      if (response.statusCode == 200) {
        return List<Map<String, dynamic>>.from(json.decode(response.body));
      } else {
        throw Exception(
          'Failed to load spot trade history: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Network error: $e');
    }
  }
}
