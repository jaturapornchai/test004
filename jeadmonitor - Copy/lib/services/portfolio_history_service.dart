import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/daily_portfolio.dart';

class PortfolioHistoryService {
  static const String _portfolioHistoryKey = 'portfolio_history';
  static const String _lastSavedDateKey = 'last_saved_date';

  // Save today's portfolio data
  Future<void> saveTodayPortfolio(DailyPortfolio portfolio) async {
    final prefs = await SharedPreferences.getInstance();

    // Get existing history
    final List<DailyPortfolio> history = await getPortfolioHistory();

    // Remove today's data if it already exists (to update)
    history.removeWhere((p) => p.date == portfolio.date);

    // Add today's data
    history.add(portfolio);

    // Keep only last 30 days
    if (history.length > 30) {
      history.sort((a, b) => a.date.compareTo(b.date));
      history.removeRange(0, history.length - 30);
    }

    // Save back to storage
    final List<String> jsonList = history
        .map((p) => jsonEncode(p.toJson()))
        .toList();
    await prefs.setStringList(_portfolioHistoryKey, jsonList);
    await prefs.setString(_lastSavedDateKey, portfolio.date);
  }

  // Get portfolio history
  Future<List<DailyPortfolio>> getPortfolioHistory() async {
    final prefs = await SharedPreferences.getInstance();
    final List<String>? jsonList = prefs.getStringList(_portfolioHistoryKey);

    if (jsonList == null) return [];

    return jsonList.map((jsonStr) {
      final Map<String, dynamic> json = jsonDecode(jsonStr);
      return DailyPortfolio.fromJson(json);
    }).toList();
  }

  // Get yesterday's portfolio for comparison
  Future<DailyPortfolio?> getYesterdayPortfolio() async {
    final history = await getPortfolioHistory();
    if (history.isEmpty) return null;

    final now = DateTime.now();
    final yesterday = DateTime(now.year, now.month, now.day - 1);
    final yesterdayStr = _formatDate(yesterday);

    try {
      return history.firstWhere((p) => p.date == yesterdayStr);
    } catch (e) {
      // If no exact yesterday data, get the most recent data
      if (history.isNotEmpty) {
        history.sort((a, b) => b.date.compareTo(a.date));
        return history.first;
      }
      return null;
    }
  }

  // Get today's date in YYYY-MM-DD format
  String getTodayDate() {
    return _formatDate(DateTime.now());
  }

  // Check if we've already saved data today
  Future<bool> hasDataForToday() async {
    final prefs = await SharedPreferences.getInstance();
    final lastSavedDate = prefs.getString(_lastSavedDateKey);
    return lastSavedDate == getTodayDate();
  }

  String _formatDate(DateTime date) {
    return '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}';
  }

  // Calculate daily P&L compared to yesterday
  Future<Map<String, double>> calculateDailyPnL(DailyPortfolio today) async {
    final yesterday = await getYesterdayPortfolio();
    if (yesterday == null) {
      return {
        'totalPnL': 0.0,
        'totalPnLPercent': 0.0,
        'binanceSpotPnL': 0.0,
        'binanceFuturesPnL': 0.0,
        'gateSpotPnL': 0.0,
        'gateFuturesPnL': 0.0,
      };
    }

    final totalPnL = today.totalValue - yesterday.totalValue;
    final totalPnLPercent = yesterday.totalValue > 0
        ? (totalPnL / yesterday.totalValue) * 100
        : 0.0;

    return {
      'totalPnL': totalPnL,
      'totalPnLPercent': totalPnLPercent,
      'binanceSpotPnL': today.binanceSpot - yesterday.binanceSpot,
      'binanceFuturesPnL': today.binanceFutures - yesterday.binanceFutures,
      'gateSpotPnL': today.gateSpot - yesterday.gateSpot,
      'gateFuturesPnL': today.gateFutures - yesterday.gateFutures,
    };
  }
}
