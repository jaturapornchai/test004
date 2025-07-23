# Crypto Portfolio Monitor

A beautiful Flutter application for real-time cryptocurrency portfolio monitoring from Binance and Gate.io exchanges.

## Features

- 🔄 **Auto-refresh every 5 seconds** - Real-time portfolio updates
- 📱 **Landscape-only mode** - Optimized for horizontal viewing
- 🎨 **Beautiful gradient UI** - Modern dark theme with colorful gradients
- 💰 **Multi-exchange support** - Binance and Gate.io integration
- 📊 **Real-time USD values** - Live price calculations
- 🚀 **Fullscreen experience** - Immersive viewing without system UI

## Screenshots

The app displays your crypto assets in a beautiful 4-column grid layout with:
- Total portfolio value in USD
- Individual asset balances and values
- Exchange identification (Binance/Gate.io)
- Auto-refreshing data every 5 seconds

## Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd jeadmonitor
   ```

2. **Install dependencies**
   ```bash
   flutter pub get
   ```

3. **Configure API keys**
   Update the API keys in `lib/providers/portfolio_provider.dart` or set environment variables:
   ```
   BINANCE_API_KEY=your_binance_api_key
   BINANCE_API_SECRET=your_binance_api_secret
   GATE_API_KEY=your_gate_api_key
   GATE_API_SECRET=your_gate_api_secret
   ```

4. **Run on Android device**
   ```bash
   flutter run
   ```

## API Configuration

The app uses the following APIs:
- **Binance API v3** - Account info and price data
- **Gate.io API v4** - Spot accounts and ticker data

Both APIs require proper authentication using HMAC signatures.

## Build Requirements

- Flutter SDK 3.8.1+
- Android SDK for device deployment
- Internet connectivity for API calls

## Usage

1. Launch the app on your Android device
2. The app will automatically:
   - Force landscape orientation
   - Enter fullscreen mode
   - Start fetching portfolio data
   - Refresh every 5 seconds

3. View your assets organized by:
   - Symbol and exchange
   - USD value
   - Total balance
   - Exchange badge (B for Binance, G for Gate.io)

## Development

The project structure:
```
lib/
├── models/         # Data models (Asset)
├── providers/      # State management (PortfolioProvider)
├── screens/        # UI screens (PortfolioScreen)
├── services/       # API clients (Binance, Gate.io)
└── main.dart      # App entry point
```

## Notes

- App is designed for Android tablets in landscape mode
- No scrolling required - all assets fit in the grid
- Automatic error handling and retry functionality
- Performance optimized for real-time updates
