import 'dart:convert';
import 'dart:io';

class ExchangeRateService {
  static const String _apiUrl =
      'https://api.exchangerate-api.com/v4/latest/USD';

  Future<double> getUsdToThbRate() async {
    try {
      final client = HttpClient();
      client.badCertificateCallback = (cert, host, port) => true;

      final request = await client.getUrl(Uri.parse(_apiUrl));
      request.headers.set('User-Agent', 'Mozilla/5.0');

      final response = await request.close();

      if (response.statusCode == 200) {
        final responseBody = await response.transform(utf8.decoder).join();
        final data = json.decode(responseBody);

        // Get THB rate from the response
        final thbRate = data['rates']['THB'];
        return double.tryParse(thbRate.toString()) ??
            35.0; // Default to 35 if fails
      } else {
        print('Failed to get exchange rate: ${response.statusCode}');
        return 35.0; // Default rate
      }
    } catch (e) {
      print('Error fetching exchange rate: $e');
      return 35.0; // Default rate in case of error
    }
  }
}
