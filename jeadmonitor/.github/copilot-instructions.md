<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Crypto Portfolio Monitor Instructions

This is a Flutter application for monitoring cryptocurrency portfolios from Binance and Gate.io exchanges.

## Key Features:

- Landscape-only orientation (no portrait mode)
- Fullscreen immersive mode (no system UI)
- Auto-refresh every 5 seconds
- Real-time portfolio tracking
- Beautiful gradient UI with dark theme

## Architecture:

- Uses Provider pattern for state management
- API clients for Binance and Gate.io
- Custom models for asset representation
- Responsive grid layout for asset display

## API Integration:

- Binance REST API v3
- Gate.io REST API v4
- HMAC-SHA256/SHA512 authentication
- Environment-based API key management

## UI/UX Guidelines:

- Dark theme with blue/orange/green gradients
- Grid layout showing 4 columns of assets
- Real-time USD value calculations
- Exchange identification badges
- Loading states and error handling

## Development Notes:

- Target Android devices only
- No scrolling required (fixed grid layout)
- Optimized for tablet landscape viewing
- Performance-focused with efficient API calls
