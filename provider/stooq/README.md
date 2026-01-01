# Stooq Provider

This directory contains the implementation for the [Stooq](https://stooq.com) data provider. Stooq is a valuable alternative for market data, but it comes with specific constraints and behaviors detailed below.

## Limitations

### Data Availability (Dividends & Yields)
Stooq's standard CSV export does not provide dividend rate or yield information. Consequently, the **Dividend** and **Yield** columns in Mop will always appear as `-` when using this provider.

### 52-Week High/Low Performance
Unlike some API-based providers that return 52-week range data in a single quote response, Stooq requires a separate historical data query to calculate these values.
- **Why it's slow**: To respect Stooq's server load and strict rate limiting, Mop fetches this historical data for one stock at a time in the background.
- **Behavior**: When you first start Mop, 52-week High and Low values will be missing (`-`). They will populate one by one as the background fetcher processes your ticker list.
- **Caching**: These values are cached locally in `~/.mop_history.json` (valid for 24 hours), so subsequent runs will display this data immediately.

### Rate Limits
Stooq strictly enforces rate limits.
- **Throttle**: The provider enforces a minimum delay (default ~5 seconds) between requests to prevent your IP from being banned.
- **Impact**: If you have a large list of tickers, the initial population of historical data (52-week high/low) will take `Number_of_Stocks * 5 seconds`.
- **Recommendation**: Keep your ticker list concise to ensure timely updates.

## Configuration
To use Stooq, ensure your `~/.moprc` file includes:
```json
"Provider": "stooq"
```
