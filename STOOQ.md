# Stooq Data Provider

This document explains the behavior and limitations of the Stooq data provider in Mop.

## Usage
To use Stooq as your data provider, update your `~/.moprc` file:

```json
{
    "Provider": "stooq",
    ...
}
```

## Limitations

### Missing Data (Dividends & Yields)
The Stooq CSV interface does not provide dividend rates or yields. These columns will show as `-` when using this provider.

### Rate Limiting & Throttling
Stooq strictly enforces rate limits. To avoid IP bans, Mop implements a mandatory ~5-second delay between requests to Stooq.

### 52-Week High/Low (Delayed Loading)
Standard quotes from Stooq do not include 52-week statistics. Mop retrieves this information by fetching historical data for each ticker individually in the background.

- **Initialization**: When Mop starts, 52-week values will be empty (`-`).
- **Progressive Loading**: The background worker fetches one stock every 5 seconds. You will see these values populate one by one in the UI.
- **Caching**: Historical data is cached in `~/.mop_history.json` for 24 hours. Once cached, the values will be available immediately on the next launch.

## Recommendations
Because of the per-ticker background fetching and strict throttling, the Stooq provider works best with a concise list of tickers. If you have dozens of stocks, it may take several minutes for the 52-week High/Low data to fully populate.
